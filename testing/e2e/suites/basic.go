package suites

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/Azure/aks-app-routing-operator/testing/e2e/infra"
	"github.com/Azure/aks-app-routing-operator/testing/e2e/logger"
	"github.com/Azure/aks-app-routing-operator/testing/e2e/manifests"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"golang.org/x/sync/errgroup"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	// basicNs is a map of namespaces that are used by the basic suite.
	// the key is the dns zone name, and the value is the namespace that
	// is used for the tests for that dns zone. Using shared namespaces
	// allow us to appropriately test upgrade scenarios.
	basicNs = make(map[string]*corev1.Namespace)
	// nsMutex is a mutex that is used to protect the basicNs map. Without this we chance concurrent goroutine map access panics
	nsMutex = sync.RWMutex{}
)

func basicSuite(in infra.Provisioned) []test {
	return []test{
		{
			name: "basic ingress",
			cfgs: builderFromInfra(in).
				withOsm(in, false, true).
				withVersions(manifests.AllUsedOperatorVersions...).
				withZones(manifests.NonZeroDnsZoneCounts, manifests.NonZeroDnsZoneCounts).
				build(),
			run: func(ctx context.Context, config *rest.Config, operator manifests.OperatorConfig) error {
				if err := clientServerTest(ctx, config, operator, uniqueNamespaceNamespacer{namespaces: basicNs}, in, nil, nil, getZoners); err != nil {
					return err
				}

				return nil
			},
		},
		{
			name: "basic service",
			cfgs: builderFromInfra(in).
				withOsm(in, false, true).
				withVersions(manifests.AllUsedOperatorVersions...).
				withZones(manifests.NonZeroDnsZoneCounts, manifests.NonZeroDnsZoneCounts).
				build(),
			run: func(ctx context.Context, config *rest.Config, operator manifests.OperatorConfig) error {
				if err := clientServerTest(ctx, config, operator, uniqueNamespaceNamespacer{namespaces: basicNs}, in, func(ingress *netv1.Ingress, service *corev1.Service, z zoner) error {
					ingress = nil
					annotations := service.GetAnnotations()
					annotations["kubernetes.azure.com/ingress-host"] = z.GetName()
					annotations["kubernetes.azure.com/tls-cert-keyvault-uri"] = z.GetCertId()
					service.SetAnnotations(annotations)

					return nil
				}, nil, getZoners); err != nil {
					return err
				}

				return nil
			},
		},
	}
}

// modifier is a function that can be used to modify the ingress and service
type modifier func(ingress *netv1.Ingress, service *corev1.Service, z zoner) error

// namespacer returns the namespace that should be used
type namespacer interface {
	getNamespace(ctx context.Context, cl client.Client, key string) (*corev1.Namespace, error)
}

type uniqueNamespaceNamespacer struct {
	namespaces map[string]*corev1.Namespace
}

func (u uniqueNamespaceNamespacer) getNamespace(ctx context.Context, cl client.Client, key string) (*corev1.Namespace, error) {
	// multiple goroutines access the same map at the same time which is not safe
	nsMutex.Lock()

	if u.namespaces == nil {
		u.namespaces = make(map[string]*corev1.Namespace)
	}

	if val, ok := u.namespaces[key]; !ok || val == nil {
		u.namespaces[key] = manifests.UncollisionedNs()
	}
	ns := u.namespaces[key]
	nsMutex.Unlock()

	if err := upsert(ctx, cl, ns); err != nil {
		return nil, fmt.Errorf("upserting ns: %w", err)
	}

	return ns, nil
}

func getZoners(ctx context.Context, c client.Client, namespacer namespacer, operator manifests.OperatorConfig, infra infra.Provisioned, serviceName *string) ([]zoner, error) {
	var zoners []zoner
	switch operator.Zones.Public {
	case manifests.DnsZoneCountNone:
	case manifests.DnsZoneCountOne:
		zs, err := toZoners(ctx, c, namespacer, infra.Zones[0])
		if err != nil {
			return nil, fmt.Errorf("converting to zoners: %w", err)
		}
		zoners = append(zoners, zs...)
	case manifests.DnsZoneCountMultiple:
		for _, z := range infra.Zones {
			zs, err := toZoners(ctx, c, namespacer, z)
			if err != nil {
				return nil, fmt.Errorf("converting to zoners: %w", err)
			}
			zoners = append(zoners, zs...)
		}
	}
	switch operator.Zones.Private {
	case manifests.DnsZoneCountNone:
	case manifests.DnsZoneCountOne:
		zs, err := toPrivateZoners(ctx, c, namespacer, infra.PrivateZones[0], infra.Cluster.GetDnsServiceIp())
		if err != nil {
			return nil, fmt.Errorf("converting to zoners: %w", err)
		}
		zoners = append(zoners, zs...)
	case manifests.DnsZoneCountMultiple:
		for _, z := range infra.PrivateZones {
			zs, err := toPrivateZoners(ctx, c, namespacer, z, infra.Cluster.GetDnsServiceIp())
			if err != nil {
				return nil, fmt.Errorf("converting to zoners: %w", err)
			}
			zoners = append(zoners, zs...)
		}
	}

	if operator.Zones.Public == manifests.DnsZoneCountNone && operator.Zones.Private == manifests.DnsZoneCountNone {
		zoners = append(zoners, zone{
			name:       fmt.Sprintf("%s.app-routing-system.svc.cluster.local:80", *serviceName),
			nameserver: infra.Cluster.GetDnsServiceIp(),
			host:       fmt.Sprintf("%s.app-routing-system.svc.cluster.local", *serviceName),
		})
	}

	return zoners, nil
}

// clientServerTest is a test that deploys a client and server application and ensures the client can reach the server.
// This is the standard test used to check traffic flow is working.
var clientServerTest = func(ctx context.Context, config *rest.Config, operator manifests.OperatorConfig, namespacer namespacer, infra infra.Provisioned, mod modifier, serviceName *string,
	getZoners func(ctx context.Context, c client.Client, namespacer namespacer, operator manifests.OperatorConfig, infra infra.Provisioned, serviceName *string) ([]zoner, error),
) error {
	lgr := logger.FromContext(ctx)
	lgr.Info("starting test")

	if serviceName == nil {
		serviceName = to.Ptr("nginx")
	}

	c, err := client.New(config, client.Options{})
	if err != nil {
		return fmt.Errorf("creating client: %w", err)
	}

	var eg errgroup.Group
	zoners, err := getZoners(ctx, c, namespacer, operator, infra, serviceName)
	for _, zone := range zoners {
		zone := zone
		eg.Go(func() error {
			lgr := logger.FromContext(ctx).With("zone", zone.GetName())
			ctx := logger.WithContext(ctx, lgr)

			ns, err := namespacer.getNamespace(ctx, c, zone.GetName())
			if err != nil {
				return fmt.Errorf("getting namespace: %w", err)
			}

			lgr = lgr.With("namespace", ns.Name)
			ctx = logger.WithContext(ctx, lgr)

			testingResources := manifests.ClientAndServer(ns.Name, zone.GetName()[:40], zone.GetNameserver(), zone.GetCertId(), zone.GetHost(), zone.GetTlsHost())
			if mod != nil {
				if err := mod(testingResources.Ingress, testingResources.Service, zone); err != nil {
					return fmt.Errorf("modifying ingress and service: %w", err)
				}
			}
			for _, object := range testingResources.Objects() {
				if err := upsert(ctx, c, object); err != nil {
					return fmt.Errorf("upserting resource: %w", err)
				}
			}

			ctx = logger.WithContext(ctx, lgr.With("client", testingResources.Client.GetName(), "clientNamespace", testingResources.Client.GetNamespace()))
			if err := waitForAvailable(ctx, c, *testingResources.Client); err != nil {
				return fmt.Errorf("waiting for client deployment to be available: %w", err)
			}

			lgr.Info("finished testing zone")
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return fmt.Errorf("testing all zones: %w", err)
	}

	lgr.Info("finished successfully")
	return nil
}

func toZoners(ctx context.Context, cl client.Client, namespacer namespacer, z infra.WithCert[infra.Zone]) ([]zoner, error) {
	name := z.Zone.GetName()
	nameserver := z.Zone.GetNameservers()[0]
	certName := z.Cert.GetName()
	certId := z.Cert.GetId()
	ns, err := namespacer.getNamespace(ctx, cl, name)
	if err != nil {
		return nil, fmt.Errorf("getting namespaces: %w", err)
	}

	return []zoner{
		zone{
			name:       name,
			nameserver: nameserver,
			certName:   certName,
			certId:     certId,
			host:       strings.ToLower(ns.Name) + "." + strings.TrimRight(name, "."),
			tlsHost:    strings.ToLower(ns.Name) + "." + strings.TrimRight(name, "."),
		},
		zone{
			name:       "wildcard" + name,
			nameserver: nameserver,
			certName:   certName,
			certId:     certId,
			host:       "wildcard." + strings.ToLower(ns.Name) + "." + strings.TrimRight(name, "."),
			tlsHost:    "*." + strings.ToLower(ns.Name) + "." + strings.TrimRight(name, "."),
		},
	}, nil
}

func toPrivateZoners(ctx context.Context, cl client.Client, namespacer namespacer, z infra.WithCert[infra.PrivateZone], nameserver string) ([]zoner, error) {
	name := z.Zone.GetName()
	certName := z.Cert.GetName()
	certId := z.Cert.GetId()
	ns, err := namespacer.getNamespace(ctx, cl, name)
	if err != nil {
		return nil, fmt.Errorf("getting namespaces: %w", err)
	}

	return []zoner{
		zone{
			name:       name,
			nameserver: nameserver,
			certName:   certName,
			certId:     certId,
			host:       strings.ToLower(ns.Name) + "." + strings.TrimRight(name, "."),
			tlsHost:    strings.ToLower(ns.Name) + "." + strings.TrimRight(name, "."),
		},
		zone{
			name:       "wildcard" + name,
			nameserver: nameserver,
			certName:   certName,
			certId:     certId,
			host:       "wildcard." + strings.ToLower(ns.Name) + "." + strings.TrimRight(name, "."),
			tlsHost:    "*." + strings.ToLower(ns.Name) + "." + strings.TrimRight(name, "."),
		},
	}, nil
}

// zoner represents a DNS endpoint and the host, nameserver, and cert information used to connect to it
type zoner interface {
	GetName() string
	GetNameserver() string
	GetCertName() string
	GetCertId() string
	GetHost() string
	GetTlsHost() string
}

type zone struct {
	name       string
	nameserver string
	certName   string
	certId     string
	host       string
	tlsHost    string
}

func (z zone) GetName() string {
	return z.name
}

func (z zone) GetNameserver() string {
	return z.nameserver
}

func (z zone) GetCertName() string {
	return z.certName
}

func (z zone) GetCertId() string {
	return z.certId
}

func (z zone) GetHost() string {
	return z.host
}

func (z zone) GetTlsHost() string {
	return z.tlsHost
}
