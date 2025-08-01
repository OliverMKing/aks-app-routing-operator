package util

import (
	"context"
	"errors"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// WiSaClientIdAnnotation is the annotation key used to specify the client ID for Workload Identity on a ServiceAccount for Workload Identity
const WiSaClientIdAnnotation = "azure.workload.identity/client-id"

func GetServiceAccountWorkloadIdentityClientId(ctx context.Context, k8sclient client.Client, saName, saNamespace string) (string, error) {
	// ensure referenced serviceaccount exists
	saObj := &corev1.ServiceAccount{}
	err := k8sclient.Get(ctx, types.NamespacedName{Name: saName, Namespace: saNamespace}, saObj)

	if client.IgnoreNotFound(err) != nil {
		return "", fmt.Errorf("failed to fetch serviceaccount to verify workload identity configuration: %w", err)
	}

	// SA wasn't found, return appropriate error
	if err != nil {
		return "", NewUserError(err, fmt.Sprintf("serviceAccount %s does not exist in namespace %s", saName, saNamespace))
	}
	// check for required annotations
	if saObj.Annotations == nil || saObj.Annotations[WiSaClientIdAnnotation] == "" {
		return "", NewUserError(errors.New("user-specified service account does not contain WI annotation"), fmt.Sprintf("serviceAccount %s was specified but does not include necessary annotation for workload identity", saName))
	}

	return saObj.Annotations[WiSaClientIdAnnotation], nil
}
