package v1alpha1

import (
	"context"
	"fmt"

	"github.com/Azure/aks-app-routing-operator/api"
	"github.com/go-logr/logr"
	netv1 "k8s.io/api/networking/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func init() {
	SchemeBuilder.Register(&NginxIngressController{}, &NginxIngressControllerList{})
}

const (
	// MaxCollisions is the maximum number of collisions allowed when generating a name for a managed resource. This corresponds to the status.CollisionCount
	MaxCollisions = 5
)

// Important: Run "make crd" to regenerate code after modifying this file
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// NginxIngressControllerSpec defines the desired state of NginxIngressController
type NginxIngressControllerSpec struct {
	// IngressClassName is the name of the IngressClass that will be used for the NGINX Ingress Controller. Defaults to metadata.name if
	// not specified.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=253
	// +kubebuilder:default:=nginx.approuting.kubernetes.azure.com
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="Value is immutable"
	// +kubebuilder:validation:Pattern=`^[a-z0-9][-a-z0-9\.]*[a-z0-9]$`
	// +kubebuilder:validation:Required
	IngressClassName string `json:"ingressClassName"`

	// ControllerNamePrefix is the name to use for the managed NGINX Ingress Controller resources.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=60
	// +kubebuilder:default:=nginx
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="Value is immutable"
	// +kubebuilder:validation:Pattern=`^[a-z0-9][-a-z0-9]*[a-z0-9]$`
	// +kubebuilder:validation:Required
	ControllerNamePrefix string `json:"controllerNamePrefix"`

	// LoadBalancerAnnotations is a map of annotations to apply to the NGINX Ingress Controller's Service. Common annotations
	// will be from the Azure LoadBalancer annotations here https://cloud-provider-azure.sigs.k8s.io/topics/loadbalancer/#loadbalancer-annotations
	// +optional
	LoadBalancerAnnotations map[string]string `json:"loadBalancerAnnotations,omitempty"`

	// LoadBalancerSourceRanges restricts traffic to the LoadBalancer Service to the specified client IPs. This can be used along with
	// deny-all annotations to restrict access  https://cloud-provider-azure.sigs.k8s.io/topics/loadbalancer/#loadbalancer-annotations
	LoadBalancerSourceRanges []string `json:"loadBalancerSourceRanges,omitempty"`

	// DefaultSSLCertificate defines whether the NginxIngressController should use a certain SSL certificate by default.
	// If this field is omitted, no default certificate will be used.
	// +optional
	DefaultSSLCertificate *DefaultSSLCertificate `json:"defaultSSLCertificate,omitempty"`

	// DefaultBackendService defines the service that the NginxIngressController should default to when given HTTP traffic with not matching known server names.
	// The controller directs traffic to the first port of the service.
	// +optional
	DefaultBackendService *NICNamespacedName `json:"defaultBackendService,omitempty"`

	// CustomHTTPErrors defines the error codes that the NginxIngressController should send to its default-backend in case of error.
	// +optional
	CustomHTTPErrors []int32 `json:"customHTTPErrors,omitempty"`

	// Scaling defines configuration options for how the Ingress Controller scales
	// +optional
	Scaling *Scaling `json:"scaling,omitempty"`

	// HTTPDisabled is a flag that disables HTTP traffic to the NginxIngressController
	// +optional
	HTTPDisabled bool `json:"httpDisabled,omitempty"`

	// LogFormat is the log format used by the Nginx Ingress Controller. See https://kubernetes.github.io/ingress-nginx/user-guide/nginx-configuration/configmap/#log-format-upstream
	// +optional
	LogFormat *string `json:"logFormat,omitempty"`

	// EnableSSLPassthrough is a flag that enables SSL passthrough for the NginxIngressController. This allows the controller to pass through SSL traffic without terminating it.
	// +optional
	EnableSSLPassthrough bool `json:"enableSSLPassthrough,omitempty"`
}

// DefaultSSLCertificate holds a secret in the form of a secret struct with name and namespace properties or a key vault uri
// +kubebuilder:validation:MaxProperties=2
// +kubebuilder:validation:XValidation:rule="(isURL(self.keyVaultURI) || !has(self.keyVaultURI))"
// +kubebuilder:validation:XValidation:rule="((self.forceSSLRedirect == true) && (has(self.secret) || has(self.keyVaultURI)) || (self.forceSSLRedirect == false))"
type DefaultSSLCertificate struct {
	// Secret is a struct that holds the name and namespace fields used for the default ssl secret
	// +optional
	Secret *Secret `json:"secret,omitempty"`

	// Secret in the form of a Key Vault URI
	// +optional
	KeyVaultURI *string `json:"keyVaultURI,omitempty"`

	// ForceSSLRedirect is a flag that sets the global value of redirects to HTTPS if there is a defined DefaultSSLCertificate
	// +kubebuilder:default:=false
	ForceSSLRedirect bool `json:"forceSSLRedirect,omitempty"`
	// forceSSLRedirect is set to false by default and will add the "forceSSLRedirect: false" property even if the user doesn't specify it.
	// If a user adds both a keyvault uri and secret the property count will be 3 since forceSSLRedirect still automatically gets added thus failing the check.
}

// Secret is a struct that holds a name and namespace to be used in DefaultSSLCertificate
type Secret struct {
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=253
	// +kubebuilder:validation:Pattern=`^[a-z0-9][-a-z0-9\.]*[a-z0-9]$`
	Name string `json:"name"`

	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=253
	// +kubebuilder:validation:Pattern=`^[a-z0-9][-a-z0-9\.]*[a-z0-9]$`
	Namespace string `json:"namespace"`
}

// NICNamespacedName is a struct that holds a name and namespace with length checking on the crd for fields other than DefaultSSLCertificate in the spec
type NICNamespacedName struct {
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=253
	// +kubebuilder:validation:Pattern=`^[a-z0-9][-a-z0-9\.]*[a-z0-9]$`
	Name string `json:"name"`

	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=253
	// +kubebuilder:validation:Pattern=`^[a-z0-9][-a-z0-9\.]*[a-z0-9]$`
	Namespace string `json:"namespace"`
}

// Scaling holds specification for how the Ingress Controller scales
// +kubebuilder:validation:XValidation:rule="(!has(self.minReplicas)) || (!has(self.maxReplicas)) || (self.minReplicas <= self.maxReplicas)"
type Scaling struct {
	// MinReplicas is the lower limit for the number of Ingress Controller replicas. It defaults to 2 pods.
	// +kubebuilder:validation:Minimum=1
	// +optional
	MinReplicas *int32 `json:"minReplicas,omitempty"`
	// MaxReplicas is the upper limit for the number of Ingress Controller replicas. It defaults to 100 pods.
	// +kubebuilder:validation:Minimum=1
	// +optional
	MaxReplicas *int32 `json:"maxReplicas,omitempty"`

	// Threshold defines how quickly the Ingress Controller pods should scale based on workload. Rapid means the Ingress Controller
	// will scale quickly and aggressively, which is the best choice for handling sudden and significant traffic spikes. Steady
	// is the opposite, prioritizing cost-effectiveness. Steady is the best choice when fewer replicas handling more work is desired or when
	// traffic isn't expected to fluctuate. Balanced is a good mix between the two that works for most use-cases. If unspecified, this field
	// defaults to balanced.
	// +kubebuilder:validation:Enum=rapid;balanced;steady;
	// +optional
	Threshold *Threshold `json:"threshold,omitempty"`
}

type Threshold string

const (
	RapidThreshold    Threshold = "rapid"
	BalancedThreshold Threshold = "balanced"
	SteadyThreshold   Threshold = "steady"
)

// NginxIngressControllerStatus defines the observed state of NginxIngressController
type NginxIngressControllerStatus struct {
	// Conditions is an array of current observed conditions for the NGINX Ingress Controller
	// +optional
	// +patchMergeKey=type
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=type
	Conditions []metav1.Condition `json:"conditions"`

	// ControllerReplicas is the desired number of replicas of the NGINX Ingress Controller
	// +optional
	ControllerReplicas int32 `json:"controllerReplicas"`

	// ControllerReadyReplicas is the number of ready replicas of the NGINX Ingress Controller deployment
	// +optional
	ControllerReadyReplicas int32 `json:"controllerReadyReplicas"`

	// ControllerAvailableReplicas is the number of available replicas of the NGINX Ingress Controller deployment
	// +optional
	ControllerAvailableReplicas int32 `json:"controllerAvailableReplicas"`

	// ControllerUnavailableReplicas is the number of unavailable replicas of the NGINX Ingress Controller deployment
	// +optional
	ControllerUnavailableReplicas int32 `json:"controllerUnavailableReplicas"`

	// Count of hash collisions for the managed resources. The App Routing Operator uses this field
	// as a collision avoidance mechanism when it needs to create the name for the managed resources.
	// +optional
	// +kubebuilder:validation:Maximum=5
	CollisionCount int32 `json:"collisionCount"`

	// ManagedResourceRefs is a list of references to the managed resources
	// +optional
	ManagedResourceRefs []ManagedObjectReference `json:"managedResourceRefs,omitempty"`
}

const (
	// ConditionTypeAvailable indicates whether the NGINX Ingress Controller is available. Its condition status is one of
	// - "True" when the NGINX Ingress Controller is available and can be used
	// - "False" when the NGINX Ingress Controller is not available and cannot offer full functionality
	// - "Unknown" when the NGINX Ingress Controller's availability cannot be determined
	ConditionTypeAvailable = "Available"

	// ConditionTypeIngressClassReady indicates whether the IngressClass exists. Its condition status is one of
	// - "True" when the IngressClass exists
	// - "False" when the IngressClass does not exist
	// - "Collision" when the IngressClass exists, but it's not owned by the NginxIngressController.
	// - "Unknown" when the IngressClass's existence cannot be determined
	ConditionTypeIngressClassReady = "IngressClassReady"

	// ConditionTypeControllerAvailable indicates whether the NGINX Ingress Controller deployment is available. Its condition status is one of
	// - "True" when the NGINX Ingress Controller deployment is available
	// - "False" when the NGINX Ingress Controller deployment is not available
	// - "Unknown" when the NGINX Ingress Controller deployment's availability cannot be determined
	ConditionTypeControllerAvailable = "ControllerAvailable"

	// ConditionTypeProgressing indicates whether the NGINX Ingress Controller availability is progressing. Its condition status is one of
	// - "True" when the NGINX Ingress Controller availability is progressing
	// - "False" when the NGINX Ingress Controller availability is not progressing
	// - "Unknown" when the NGINX Ingress Controller availability's progress cannot be determined
	ConditionTypeProgressing = "Progressing"
)

// ManagedObjectReference is a reference to an object
type ManagedObjectReference struct {
	// Name is the name of the managed object
	Name string `json:"name"`

	// Namespace is the namespace of the managed object. If not specified, the resource is cluster-scoped
	// +optional
	Namespace string `json:"namespace"`

	// Kind is the kind of the managed object
	Kind string `json:"kind"`

	// APIGroup is the API group of the managed object. If not specified, the resource is in the core API group
	// +optional
	APIGroup string `json:"apiGroup"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster,shortName=nic
//+kubebuilder:printcolumn:name="IngressClass",type="string",JSONPath=`.spec.ingressClassName`
//+kubebuilder:printcolumn:name="ControllerNamePrefix",type="string",JSONPath=`.spec.controllerNamePrefix`
//+kubebuilder:printcolumn:name="Available",type="string",JSONPath=`.status.conditions[?(@.type=="Available")].status`

// NginxIngressController is the Schema for the nginxingresscontrollers API
type NginxIngressController struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// +required
	// +kubebuilder:default:={"ingressClassName":"nginx.approuting.kubernetes.azure.com","controllerNamePrefix":"nginx"}
	Spec NginxIngressControllerSpec `json:"spec"` // ^ for the above thing https://github.com/kubernetes-sigs/controller-tools/issues/622 defaulting doesn't cascade, so we have to define it all. Comment on this line so it's not in crd spec.

	// +optional
	Status NginxIngressControllerStatus `json:"status,omitempty"`
}

func (n *NginxIngressController) GetGeneration() int64 {
	return n.Generation
}

func (n *NginxIngressController) GetConditions() *[]metav1.Condition {
	return &n.Status.Conditions
}

func (n *NginxIngressController) GetCondition(t string) *metav1.Condition {
	return meta.FindStatusCondition(n.Status.Conditions, t)
}

func (n *NginxIngressController) SetCondition(c metav1.Condition) {
	api.VerifyAndSetCondition(n, c)
}

// Collides returns whether the fields in this NginxIngressController would collide with an existing resources making it
// impossible for this NginxIngressController to become available. This should be run before an NginxIngressController is created.
// Returns whether there's a collision, the collision reason, and an error if one occurred. The collision reason is something that
// the user can use to understand and resolve.
func (n *NginxIngressController) Collides(ctx context.Context, cl client.Client) (bool, string, error) {
	lgr := logr.FromContextOrDiscard(ctx).WithValues("name", n.Name, "ingressClassName", n.Spec.IngressClassName)
	lgr.Info("checking for NginxIngressController collisions")

	// check for NginxIngressController collisions
	lgr.Info("checking for NginxIngressController collision")
	var nginxIngressControllerList NginxIngressControllerList
	if err := cl.List(ctx, &nginxIngressControllerList); err != nil {
		lgr.Error(err, "listing NginxIngressControllers")
		return false, "", fmt.Errorf("listing NginxIngressControllers: %w", err)
	}

	for _, nic := range nginxIngressControllerList.Items {
		if nic.Spec.IngressClassName == n.Spec.IngressClassName && nic.Name != n.Name {
			lgr.Info("NginxIngressController collision found")
			return true, fmt.Sprintf("spec.ingressClassName \"%s\" is invalid because NginxIngressController \"%s\" already uses IngressClass \"%[1]s\"", n.Spec.IngressClassName, nic.Name), nil
		}
	}

	// Check for an IngressClass collision.
	// This is purposefully after the NginxIngressController check because if the collision is through an NginxIngressController
	// that's the one we want to report as the reason since the user action for fixing that would involve working with the NginxIngressController
	// resource rather than the IngressClass resource.
	lgr.Info("checking for IngressClass collision")
	ic := &netv1.IngressClass{
		ObjectMeta: metav1.ObjectMeta{
			Name: n.Spec.IngressClassName,
		},
	}
	err := cl.Get(ctx, types.NamespacedName{Name: ic.Name}, ic)
	if err == nil {
		lgr.Info("IngressClass collision found")
		return true, fmt.Sprintf("spec.ingressClassName \"%s\" is invalid because IngressClass \"%[1]s\" already exists", n.Spec.IngressClassName), nil
	}
	if !k8serrors.IsNotFound(err) {
		lgr.Error(err, "checking for IngressClass collisions")
		return false, "", fmt.Errorf("checking for IngressClass collisions: %w", err)
	}

	return false, "", nil
}

//+kubebuilder:object:root=true
//+kubebuilder:resource:scope=Cluster

// NginxIngressControllerList contains a list of NginxIngressController
type NginxIngressControllerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NginxIngressController `json:"items"`
}
