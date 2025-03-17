//go:build !ignore_autogenerated

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterExternalDNS) DeepCopyInto(out *ClusterExternalDNS) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterExternalDNS.
func (in *ClusterExternalDNS) DeepCopy() *ClusterExternalDNS {
	if in == nil {
		return nil
	}
	out := new(ClusterExternalDNS)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ClusterExternalDNS) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterExternalDNSList) DeepCopyInto(out *ClusterExternalDNSList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ClusterExternalDNS, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterExternalDNSList.
func (in *ClusterExternalDNSList) DeepCopy() *ClusterExternalDNSList {
	if in == nil {
		return nil
	}
	out := new(ClusterExternalDNSList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ClusterExternalDNSList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterExternalDNSSpec) DeepCopyInto(out *ClusterExternalDNSSpec) {
	*out = *in
	if in.DNSZoneResourceIDs != nil {
		in, out := &in.DNSZoneResourceIDs, &out.DNSZoneResourceIDs
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.ResourceTypes != nil {
		in, out := &in.ResourceTypes, &out.ResourceTypes
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	out.Identity = in.Identity
	if in.Filters != nil {
		in, out := &in.Filters, &out.Filters
		*out = new(ExternalDNSFilters)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterExternalDNSSpec.
func (in *ClusterExternalDNSSpec) DeepCopy() *ClusterExternalDNSSpec {
	if in == nil {
		return nil
	}
	out := new(ClusterExternalDNSSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterExternalDNSStatus) DeepCopyInto(out *ClusterExternalDNSStatus) {
	*out = *in
	in.ExternalDNSStatus.DeepCopyInto(&out.ExternalDNSStatus)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterExternalDNSStatus.
func (in *ClusterExternalDNSStatus) DeepCopy() *ClusterExternalDNSStatus {
	if in == nil {
		return nil
	}
	out := new(ClusterExternalDNSStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DefaultSSLCertificate) DeepCopyInto(out *DefaultSSLCertificate) {
	*out = *in
	if in.Secret != nil {
		in, out := &in.Secret, &out.Secret
		*out = new(Secret)
		**out = **in
	}
	if in.KeyVaultURI != nil {
		in, out := &in.KeyVaultURI, &out.KeyVaultURI
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DefaultSSLCertificate.
func (in *DefaultSSLCertificate) DeepCopy() *DefaultSSLCertificate {
	if in == nil {
		return nil
	}
	out := new(DefaultSSLCertificate)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ExternalDNS) DeepCopyInto(out *ExternalDNS) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ExternalDNS.
func (in *ExternalDNS) DeepCopy() *ExternalDNS {
	if in == nil {
		return nil
	}
	out := new(ExternalDNS)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ExternalDNS) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ExternalDNSFilters) DeepCopyInto(out *ExternalDNSFilters) {
	*out = *in
	if in.GatewayLabelSelector != nil {
		in, out := &in.GatewayLabelSelector, &out.GatewayLabelSelector
		*out = new(string)
		**out = **in
	}
	if in.RouteAndIngressLabelSelector != nil {
		in, out := &in.RouteAndIngressLabelSelector, &out.RouteAndIngressLabelSelector
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ExternalDNSFilters.
func (in *ExternalDNSFilters) DeepCopy() *ExternalDNSFilters {
	if in == nil {
		return nil
	}
	out := new(ExternalDNSFilters)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ExternalDNSIdentity) DeepCopyInto(out *ExternalDNSIdentity) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ExternalDNSIdentity.
func (in *ExternalDNSIdentity) DeepCopy() *ExternalDNSIdentity {
	if in == nil {
		return nil
	}
	out := new(ExternalDNSIdentity)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ExternalDNSList) DeepCopyInto(out *ExternalDNSList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ExternalDNS, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ExternalDNSList.
func (in *ExternalDNSList) DeepCopy() *ExternalDNSList {
	if in == nil {
		return nil
	}
	out := new(ExternalDNSList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ExternalDNSList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ExternalDNSSpec) DeepCopyInto(out *ExternalDNSSpec) {
	*out = *in
	if in.DNSZoneResourceIDs != nil {
		in, out := &in.DNSZoneResourceIDs, &out.DNSZoneResourceIDs
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.ResourceTypes != nil {
		in, out := &in.ResourceTypes, &out.ResourceTypes
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	out.Identity = in.Identity
	if in.Filters != nil {
		in, out := &in.Filters, &out.Filters
		*out = new(ExternalDNSFilters)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ExternalDNSSpec.
func (in *ExternalDNSSpec) DeepCopy() *ExternalDNSSpec {
	if in == nil {
		return nil
	}
	out := new(ExternalDNSSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ExternalDNSStatus) DeepCopyInto(out *ExternalDNSStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]v1.Condition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.ManagedResourceRefs != nil {
		in, out := &in.ManagedResourceRefs, &out.ManagedResourceRefs
		*out = make([]ManagedObjectReference, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ExternalDNSStatus.
func (in *ExternalDNSStatus) DeepCopy() *ExternalDNSStatus {
	if in == nil {
		return nil
	}
	out := new(ExternalDNSStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ManagedCertificate) DeepCopyInto(out *ManagedCertificate) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ManagedCertificate.
func (in *ManagedCertificate) DeepCopy() *ManagedCertificate {
	if in == nil {
		return nil
	}
	out := new(ManagedCertificate)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ManagedCertificate) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ManagedCertificateList) DeepCopyInto(out *ManagedCertificateList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ManagedCertificate, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ManagedCertificateList.
func (in *ManagedCertificateList) DeepCopy() *ManagedCertificateList {
	if in == nil {
		return nil
	}
	out := new(ManagedCertificateList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ManagedCertificateList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ManagedCertificateSpec) DeepCopyInto(out *ManagedCertificateSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ManagedCertificateSpec.
func (in *ManagedCertificateSpec) DeepCopy() *ManagedCertificateSpec {
	if in == nil {
		return nil
	}
	out := new(ManagedCertificateSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ManagedCertificateStatus) DeepCopyInto(out *ManagedCertificateStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ManagedCertificateStatus.
func (in *ManagedCertificateStatus) DeepCopy() *ManagedCertificateStatus {
	if in == nil {
		return nil
	}
	out := new(ManagedCertificateStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ManagedObjectReference) DeepCopyInto(out *ManagedObjectReference) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ManagedObjectReference.
func (in *ManagedObjectReference) DeepCopy() *ManagedObjectReference {
	if in == nil {
		return nil
	}
	out := new(ManagedObjectReference)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NICNamespacedName) DeepCopyInto(out *NICNamespacedName) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NICNamespacedName.
func (in *NICNamespacedName) DeepCopy() *NICNamespacedName {
	if in == nil {
		return nil
	}
	out := new(NICNamespacedName)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NginxIngressController) DeepCopyInto(out *NginxIngressController) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NginxIngressController.
func (in *NginxIngressController) DeepCopy() *NginxIngressController {
	if in == nil {
		return nil
	}
	out := new(NginxIngressController)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *NginxIngressController) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NginxIngressControllerList) DeepCopyInto(out *NginxIngressControllerList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]NginxIngressController, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NginxIngressControllerList.
func (in *NginxIngressControllerList) DeepCopy() *NginxIngressControllerList {
	if in == nil {
		return nil
	}
	out := new(NginxIngressControllerList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *NginxIngressControllerList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NginxIngressControllerSpec) DeepCopyInto(out *NginxIngressControllerSpec) {
	*out = *in
	if in.LoadBalancerAnnotations != nil {
		in, out := &in.LoadBalancerAnnotations, &out.LoadBalancerAnnotations
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.DefaultSSLCertificate != nil {
		in, out := &in.DefaultSSLCertificate, &out.DefaultSSLCertificate
		*out = new(DefaultSSLCertificate)
		(*in).DeepCopyInto(*out)
	}
	if in.DefaultBackendService != nil {
		in, out := &in.DefaultBackendService, &out.DefaultBackendService
		*out = new(NICNamespacedName)
		**out = **in
	}
	if in.CustomHTTPErrors != nil {
		in, out := &in.CustomHTTPErrors, &out.CustomHTTPErrors
		*out = make([]int32, len(*in))
		copy(*out, *in)
	}
	if in.Scaling != nil {
		in, out := &in.Scaling, &out.Scaling
		*out = new(Scaling)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NginxIngressControllerSpec.
func (in *NginxIngressControllerSpec) DeepCopy() *NginxIngressControllerSpec {
	if in == nil {
		return nil
	}
	out := new(NginxIngressControllerSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NginxIngressControllerStatus) DeepCopyInto(out *NginxIngressControllerStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]v1.Condition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.ManagedResourceRefs != nil {
		in, out := &in.ManagedResourceRefs, &out.ManagedResourceRefs
		*out = make([]ManagedObjectReference, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NginxIngressControllerStatus.
func (in *NginxIngressControllerStatus) DeepCopy() *NginxIngressControllerStatus {
	if in == nil {
		return nil
	}
	out := new(NginxIngressControllerStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Scaling) DeepCopyInto(out *Scaling) {
	*out = *in
	if in.MinReplicas != nil {
		in, out := &in.MinReplicas, &out.MinReplicas
		*out = new(int32)
		**out = **in
	}
	if in.MaxReplicas != nil {
		in, out := &in.MaxReplicas, &out.MaxReplicas
		*out = new(int32)
		**out = **in
	}
	if in.Threshold != nil {
		in, out := &in.Threshold, &out.Threshold
		*out = new(Threshold)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Scaling.
func (in *Scaling) DeepCopy() *Scaling {
	if in == nil {
		return nil
	}
	out := new(Scaling)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Secret) DeepCopyInto(out *Secret) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Secret.
func (in *Secret) DeepCopy() *Secret {
	if in == nil {
		return nil
	}
	out := new(Secret)
	in.DeepCopyInto(out)
	return out
}
