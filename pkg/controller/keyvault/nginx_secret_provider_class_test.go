// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package keyvault

import (
	"context"
	"testing"

	"github.com/Azure/aks-app-routing-operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/go-logr/logr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	secv1 "sigs.k8s.io/secrets-store-csi-driver/apis/v1"

	"github.com/Azure/aks-app-routing-operator/pkg/config"
	"github.com/Azure/aks-app-routing-operator/pkg/controller/metrics"
	"github.com/Azure/aks-app-routing-operator/pkg/controller/testutils"
	"github.com/Azure/aks-app-routing-operator/pkg/manifests"
	"github.com/Azure/aks-app-routing-operator/pkg/util"
)

var (
	spcTestKeyVaultURI           = "https://testvault.vault.azure.net/certificates/testcert/f8982febc6894c0697b884f946fb1a34"
	spcTestNginxIngressClassName = "webapprouting.kubernetes.azure.com"
	spcTestNginxIngress          = &v1alpha1.NginxIngressController{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-nic",
			Namespace: "test-namespace",
		},
		Spec: v1alpha1.NginxIngressControllerSpec{
			IngressClassName: spcTestNginxIngressClassName,
			DefaultSSLCertificate: &v1alpha1.DefaultSSLCertificate{
				Secret:      nil,
				KeyVaultURI: &spcTestKeyVaultURI},
		},
	}
)

func TestNginxSecretProviderClassReconcilerIntegration(t *testing.T) {
	// Create the nic
	ctx := context.Background()
	ctx = logr.NewContext(ctx, logr.Discard())
	nic := spcTestNginxIngress.DeepCopy()

	scheme := runtime.NewScheme()
	require.NoError(t, v1alpha1.AddToScheme(scheme))
	require.NoError(t, secv1.AddToScheme(scheme))

	c := fake.NewClientBuilder().WithScheme(scheme).WithObjects(nic).Build()

	recorder := record.NewFakeRecorder(10)
	n := &NginxSecretProviderClassReconciler{
		client: c,
		config: &config.Config{
			TenantID:    "test-tenant-id",
			MSIClientID: "test-msi-client-id",
		},
		events: recorder,
	}

	// Create the secret provider class
	req := ctrl.Request{NamespacedName: types.NamespacedName{Namespace: nic.Namespace, Name: nic.Name}}
	beforeErrCount := testutils.GetErrMetricCount(t, nginxSecretProviderControllerName)
	beforeRequestCount := testutils.GetReconcileMetricCount(t, nginxSecretProviderControllerName, metrics.LabelSuccess)
	_, err := n.Reconcile(ctx, req)
	require.NoError(t, err)

	require.Equal(t, testutils.GetErrMetricCount(t, nginxSecretProviderControllerName), beforeErrCount)
	require.Greater(t, testutils.GetReconcileMetricCount(t, nginxSecretProviderControllerName, metrics.LabelSuccess), beforeRequestCount)

	// Prove it exists
	spc := &secv1.SecretProviderClass{}
	spc.Name = DefaultNginxCertName(nic)
	spc.Namespace = n.config.NS
	require.NoError(t, c.Get(ctx, client.ObjectKeyFromObject(spc), spc))

	expected := &secv1.SecretProviderClass{
		Spec: secv1.SecretProviderClassSpec{
			Provider: "azure",
			Parameters: map[string]string{
				"keyvaultName":           "testvault",
				"objects":                "{\"array\":[\"{\\\"objectName\\\":\\\"testcert\\\",\\\"objectType\\\":\\\"secret\\\",\\\"objectVersion\\\":\\\"f8982febc6894c0697b884f946fb1a34\\\"}\"]}",
				"tenantId":               n.config.TenantID,
				"useVMManagedIdentity":   "true",
				"userAssignedIdentityID": n.config.MSIClientID,
			},
			SecretObjects: []*secv1.SecretObject{{
				SecretName: spc.Name,
				Type:       "kubernetes.io/tls",
				Data: []*secv1.SecretObjectData{
					{ObjectName: "testcert", Key: "tls.key"},
					{ObjectName: "testcert", Key: "tls.crt"},
				},
			}},
		},
	}
	assert.Equal(t, expected.Spec, spc.Spec)

	// Check for idempotence
	beforeErrCount = testutils.GetErrMetricCount(t, nginxSecretProviderControllerName)
	beforeRequestCount = testutils.GetReconcileMetricCount(t, nginxSecretProviderControllerName, metrics.LabelSuccess)
	_, err = n.Reconcile(ctx, req)
	require.NoError(t, err)
	require.Equal(t, testutils.GetErrMetricCount(t, nginxSecretProviderControllerName), beforeErrCount)
	require.Greater(t, testutils.GetReconcileMetricCount(t, nginxSecretProviderControllerName, metrics.LabelSuccess), beforeRequestCount)

	// Remove the cert's version from the nic
	testKvUri := "https://testvault.vault.azure.net/certificates/testcert"
	nic.Spec.DefaultSSLCertificate.KeyVaultURI = &testKvUri
	require.NoError(t, n.client.Update(ctx, nic))
	beforeErrCount = testutils.GetErrMetricCount(t, nginxSecretProviderControllerName)
	beforeRequestCount = testutils.GetReconcileMetricCount(t, nginxSecretProviderControllerName, metrics.LabelSuccess)
	_, err = n.Reconcile(ctx, req)
	require.NoError(t, err)
	require.Equal(t, testutils.GetErrMetricCount(t, nginxSecretProviderControllerName), beforeErrCount)
	require.Greater(t, testutils.GetReconcileMetricCount(t, nginxSecretProviderControllerName, metrics.LabelSuccess), beforeRequestCount)

	// Prove the objectVersion property was removed
	require.NoError(t, c.Get(ctx, client.ObjectKeyFromObject(spc), spc))
	expected.Spec.Parameters["objects"] = "{\"array\":[\"{\\\"objectName\\\":\\\"testcert\\\",\\\"objectType\\\":\\\"secret\\\"}\"]}"
	assert.Equal(t, expected.Spec, spc.Spec)

	// Remove the cert annotation from the nic
	nic.Spec.DefaultSSLCertificate.KeyVaultURI = nil
	require.NoError(t, n.client.Update(ctx, nic))
	beforeErrCount = testutils.GetErrMetricCount(t, nginxSecretProviderControllerName)
	beforeRequestCount = testutils.GetReconcileMetricCount(t, nginxSecretProviderControllerName, metrics.LabelSuccess)
	_, err = n.Reconcile(ctx, req)
	require.NoError(t, err)
	require.Equal(t, testutils.GetErrMetricCount(t, nginxSecretProviderControllerName), beforeErrCount)
	require.Greater(t, testutils.GetReconcileMetricCount(t, nginxSecretProviderControllerName, metrics.LabelSuccess), beforeRequestCount)

	// Prove secret class was removed
	require.True(t, errors.IsNotFound(c.Get(ctx, client.ObjectKeyFromObject(spc), spc)))

	// Check for idempotence
	beforeErrCount = testutils.GetErrMetricCount(t, nginxSecretProviderControllerName)
	beforeRequestCount = testutils.GetReconcileMetricCount(t, nginxSecretProviderControllerName, metrics.LabelSuccess)
	_, err = n.Reconcile(ctx, req)
	require.NoError(t, err)
	require.Equal(t, testutils.GetErrMetricCount(t, nginxSecretProviderControllerName), beforeErrCount)
	require.Greater(t, testutils.GetReconcileMetricCount(t, nginxSecretProviderControllerName, metrics.LabelSuccess), beforeRequestCount)
}

func TestNginxSecretProviderClassReconcilerIntegrationWithoutSPCLabels(t *testing.T) {
	// Create the nic
	ctx := context.Background()
	ctx = logr.NewContext(ctx, logr.Discard())
	nic := spcTestNginxIngress.DeepCopy()

	scheme := runtime.NewScheme()
	require.NoError(t, v1alpha1.AddToScheme(scheme))
	require.NoError(t, secv1.AddToScheme(scheme))

	c := fake.NewClientBuilder().WithScheme(scheme).WithObjects(nic).Build()

	recorder := record.NewFakeRecorder(10)
	n := &NginxSecretProviderClassReconciler{
		client: c,
		config: &config.Config{
			TenantID:    "test-tenant-id",
			MSIClientID: "test-msi-client-id",
		},
		events: recorder,
	}

	// Create the secret provider class
	req := ctrl.Request{NamespacedName: types.NamespacedName{Namespace: nic.Namespace, Name: nic.Name}}
	beforeErrCount := testutils.GetErrMetricCount(t, nginxSecretProviderControllerName)
	beforeRequestCount := testutils.GetReconcileMetricCount(t, nginxSecretProviderControllerName, metrics.LabelSuccess)
	_, err := n.Reconcile(ctx, req)
	require.NoError(t, err)

	require.Equal(t, testutils.GetErrMetricCount(t, nginxSecretProviderControllerName), beforeErrCount)
	require.Greater(t, testutils.GetReconcileMetricCount(t, nginxSecretProviderControllerName, metrics.LabelSuccess), beforeRequestCount)

	spc := &secv1.SecretProviderClass{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "secrets-store.csi.x-k8s.io/v1",
			Kind:       "SecretProviderClass",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      DefaultNginxCertName(nic),
			Namespace: n.config.NS,
			Labels:    manifests.GetTopLevelLabels(),
			OwnerReferences: []metav1.OwnerReference{{
				APIVersion: nic.APIVersion,
				Controller: util.ToPtr(true),
				Kind:       nic.Kind,
				Name:       nic.Name,
				UID:        nic.UID,
			}},
		},
	}

	// Get secret provider class
	require.False(t, errors.IsNotFound(c.Get(ctx, client.ObjectKeyFromObject(spc), spc)))
	assert.Equal(t, len(manifests.GetTopLevelLabels()), len(spc.Labels))

	// Remove the labels
	spc.Labels = map[string]string{}
	require.NoError(t, n.client.Update(ctx, spc))
	assert.Equal(t, 0, len(spc.Labels))

	// Remove the cert uri from the nic
	nic.Spec.DefaultSSLCertificate.KeyVaultURI = nil
	require.NoError(t, n.client.Update(ctx, nic))

	// Reconcile both changes
	beforeErrCount = testutils.GetErrMetricCount(t, nginxSecretProviderControllerName)
	beforeRequestCount = testutils.GetReconcileMetricCount(t, nginxSecretProviderControllerName, metrics.LabelSuccess)
	_, err = n.Reconcile(ctx, req)
	require.NoError(t, err)
	require.Equal(t, testutils.GetErrMetricCount(t, nginxSecretProviderControllerName), beforeErrCount)
	require.Greater(t, testutils.GetReconcileMetricCount(t, nginxSecretProviderControllerName, metrics.LabelSuccess), beforeRequestCount)

	// Prove secret class was not removed
	require.False(t, errors.IsNotFound(c.Get(ctx, client.ObjectKeyFromObject(spc), spc)))
	assert.Equal(t, 0, len(spc.Labels))
	require.NoError(t, c.Get(ctx, client.ObjectKeyFromObject(spc), spc))

	// Check secret provider class Spec after Reconcile
	expected := &secv1.SecretProviderClass{
		Spec: secv1.SecretProviderClassSpec{
			Provider: "azure",
			Parameters: map[string]string{
				"keyvaultName":           "testvault",
				"objects":                "{\"array\":[\"{\\\"objectName\\\":\\\"testcert\\\",\\\"objectType\\\":\\\"secret\\\",\\\"objectVersion\\\":\\\"f8982febc6894c0697b884f946fb1a34\\\"}\"]}",
				"tenantId":               n.config.TenantID,
				"useVMManagedIdentity":   "true",
				"userAssignedIdentityID": n.config.MSIClientID,
			},
			SecretObjects: []*secv1.SecretObject{{
				SecretName: spc.Name,
				Type:       "kubernetes.io/tls",
				Data: []*secv1.SecretObjectData{
					{ObjectName: "testcert", Key: "tls.key"},
					{ObjectName: "testcert", Key: "tls.crt"},
				},
			}},
		},
	}
	assert.Equal(t, expected.Spec, spc.Spec)

	// Check for idempotence
	beforeErrCount = testutils.GetErrMetricCount(t, nginxSecretProviderControllerName)
	beforeRequestCount = testutils.GetReconcileMetricCount(t, nginxSecretProviderControllerName, metrics.LabelSuccess)
	_, err = n.Reconcile(ctx, req)
	require.NoError(t, err)
	require.Equal(t, testutils.GetErrMetricCount(t, nginxSecretProviderControllerName), beforeErrCount)
	require.Greater(t, testutils.GetReconcileMetricCount(t, nginxSecretProviderControllerName, metrics.LabelSuccess), beforeRequestCount)
}

func TestNginxSecretProviderClassReconcilerInvalidURL(t *testing.T) {
	ctx := context.Background()
	ctx = logr.NewContext(ctx, logr.Discard())
	// Create the nic
	invalidUri := "inv@lid URL"
	nic := spcTestNginxIngress.DeepCopy()
	nic.Spec.DefaultSSLCertificate.KeyVaultURI = &invalidUri

	scheme := runtime.NewScheme()
	require.NoError(t, v1alpha1.AddToScheme(scheme))
	require.NoError(t, secv1.AddToScheme(scheme))
	c := fake.NewClientBuilder().WithScheme(scheme).WithObjects(nic).Build()
	recorder := record.NewFakeRecorder(10)
	n := &NginxSecretProviderClassReconciler{
		client: c,
		config: &config.Config{
			TenantID:    "test-tenant-id",
			MSIClientID: "test-msi-client-id",
		},
		events: recorder,
	}

	metrics.InitControllerMetrics(nginxSecretProviderControllerName)

	// get the before value of the error metrics
	beforeErrCount := testutils.GetErrMetricCount(t, nginxSecretProviderControllerName)
	beforeRequestCount := testutils.GetReconcileMetricCount(t, nginxSecretProviderControllerName, metrics.LabelError)

	req := ctrl.Request{NamespacedName: types.NamespacedName{Namespace: nic.Namespace, Name: nic.Name}}
	_, err := n.Reconcile(ctx, req)
	require.NoError(t, err)

	assert.Equal(t, "Warning InvalidInput error while processing Keyvault reference: invalid secret uri: inv@lid URL", <-recorder.Events)
	//even though no error was returned, we should expect the error count to be incremented
	afterErrCount := testutils.GetErrMetricCount(t, nginxSecretProviderControllerName)
	afterRequestCount := testutils.GetReconcileMetricCount(t, nginxSecretProviderControllerName, metrics.LabelError)

	// no new errors from invalid url
	assert.Equal(t, afterErrCount, beforeErrCount)
	assert.Equal(t, afterRequestCount, beforeRequestCount)
}
