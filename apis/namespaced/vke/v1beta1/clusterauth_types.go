package v1beta1

import (
	xpv2 "github.com/crossplane/crossplane/apis/v2/core/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// ClusterAuthParameters are the configurable fields of a ClusterAuth.
type ClusterAuthParameters struct {
	// ClusterID is the ID of the VKE cluster to fetch the kubeconfig of.
	// Either ClusterID, ClusterIDRef or ClusterIDSelector has to be given.
	// +crossplane:generate:reference:type=github.com/upbound/provider-vultr/apis/namespaced/vke/v1beta1.Kubernetes
	// +kubebuilder:validation:Optional
	ClusterID *string `json:"clusterId,omitempty"`

	// Reference to a Kubernetes in vke to populate clusterId.
	// +kubebuilder:validation:Optional
	ClusterIDRef *xpv2.NamespacedReference `json:"clusterIdRef,omitempty"`

	// Selector for a Kubernetes in vke to populate clusterId.
	// +kubebuilder:validation:Optional
	ClusterIDSelector *xpv2.NamespacedSelector `json:"clusterIdSelector,omitempty"`
}

// ClusterAuthObservation are the observable fields of a ClusterAuth.
type ClusterAuthObservation struct {
	// Endpoint is the API server address of the VKE cluster, as found in its
	// kubeconfig.
	Endpoint *string `json:"endpoint,omitempty"`
}

// ClusterAuthSpec defines the desired state of ClusterAuth
type ClusterAuthSpec struct {
	xpv2.ManagedResourceSpec `json:",inline"`
	ForProvider              ClusterAuthParameters `json:"forProvider"`
}

// ClusterAuthStatus defines the observed state of ClusterAuth.
type ClusterAuthStatus struct {
	xpv2.ManagedResourceStatus `json:",inline"`
	AtProvider                 ClusterAuthObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:storageversion

// ClusterAuth publishes the kubeconfig of a VKE cluster as connection details.
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="EXTERNAL-NAME",type="string",JSONPath=".metadata.annotations.crossplane\\.io/external-name"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:resource:scope=Namespaced,categories={crossplane,managed,vultr}
type ClusterAuth struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              ClusterAuthSpec   `json:"spec"`
	Status            ClusterAuthStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ClusterAuthList contains a list of ClusterAuths
type ClusterAuthList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ClusterAuth `json:"items"`
}

// Repository type metadata.
var (
	ClusterAuth_Kind             = "ClusterAuth"
	ClusterAuth_GroupKind        = schema.GroupKind{Group: CRDGroup, Kind: ClusterAuth_Kind}.String()
	ClusterAuth_KindAPIVersion   = ClusterAuth_Kind + "." + CRDGroupVersion.String()
	ClusterAuth_GroupVersionKind = CRDGroupVersion.WithKind(ClusterAuth_Kind)
)

func init() {
	SchemeBuilder.Register(&ClusterAuth{}, &ClusterAuthList{})
}
