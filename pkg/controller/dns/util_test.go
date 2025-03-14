package dns

import (
	"testing"

	"github.com/Azure/aks-app-routing-operator/api/v1alpha1"
	"github.com/Azure/aks-app-routing-operator/pkg/manifests"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/stretchr/testify/require"
)

type mockDnsConfig struct {
	tenantId            string
	inputServiceAccount string
	resourceNamespace   string
	inputResourceName   string
	resourceTypes       []string
	dnsZoneresourceIDs  []string
	filters             *v1alpha1.ExternalDNSFilters
	namespaced          bool
}

var sampleMockDnsConfig = mockDnsConfig{
	tenantId:            "mock-tenant-id",
	inputServiceAccount: "mock-service-account",
	resourceNamespace:   "mock-namespace",
	inputResourceName:   "mock-resource-name",
	resourceTypes:       []string{"ingress", "gateway"},
	dnsZoneresourceIDs:  []string{"mock-dns-zone-id"},
	filters: &v1alpha1.ExternalDNSFilters{
		GatewayLabelSelector:         to.Ptr("test=test"),
		RouteAndIngressLabelSelector: to.Ptr("test=othertest"),
	},
	namespaced: true,
}

func (m mockDnsConfig) GetTenantId() string {
	return m.tenantId
}

func (m mockDnsConfig) GetInputServiceAccount() string {
	return m.inputServiceAccount
}

func (m mockDnsConfig) GetResourceNamespace() string {
	return m.resourceNamespace
}

func (m mockDnsConfig) GetInputResourceName() string {
	return m.inputResourceName
}

func (m mockDnsConfig) GetResourceTypes() []string {
	return m.resourceTypes
}

func (m mockDnsConfig) GetDnsZoneresourceIDs() []string {
	return m.dnsZoneresourceIDs
}

func (m mockDnsConfig) GetFilters() *v1alpha1.ExternalDNSFilters {
	return m.filters
}

func (m mockDnsConfig) GetNamespaced() bool {
	return m.namespaced
}

func Test_buildInputDNSConfig(t *testing.T) {
	config := buildInputDNSConfig(sampleMockDnsConfig)

	require.Equal(t, config.TenantId, sampleMockDnsConfig.tenantId)
	require.Equal(t, config.InputServiceAccount, sampleMockDnsConfig.inputServiceAccount)
	require.Equal(t, config.Namespace, sampleMockDnsConfig.resourceNamespace)
	require.Equal(t, config.InputResourceName, sampleMockDnsConfig.inputResourceName)
	require.Equal(t, config.ResourceTypes, map[manifests.ResourceType]struct{}{
		manifests.ResourceTypeIngress: {},
		manifests.ResourceTypeGateway: {},
	})
	require.Equal(t, config.DnsZoneresourceIDs, sampleMockDnsConfig.dnsZoneresourceIDs)
	require.Equal(t, config.Filters, sampleMockDnsConfig.filters)
	require.Equal(t, config.IsNamespaced, sampleMockDnsConfig.namespaced)
}

func Test_extractResourceTypes(t *testing.T) {
	for _, tc := range []struct {
		rt       []string
		expected map[manifests.ResourceType]struct{}
	}{
		{
			rt: []string{"ingress", "gateway"},
			expected: map[manifests.ResourceType]struct{}{
				manifests.ResourceTypeIngress: {},
				manifests.ResourceTypeGateway: {},
			},
		},
		{
			rt: []string{"unknown", "gateway"},
			expected: map[manifests.ResourceType]struct{}{
				manifests.ResourceTypeGateway: {},
			},
		},
		{
			rt: []string{"ingress"},
			expected: map[manifests.ResourceType]struct{}{
				manifests.ResourceTypeIngress: {},
			},
		},
	} {
		result := extractResourceTypes(tc.rt)
		require.Equal(t, tc.expected, result)
	}

}
