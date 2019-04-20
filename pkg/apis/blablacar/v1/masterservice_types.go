package v1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:openapi-gen=true
type MasterServiceCallback struct {
	Port int32  `json:"port"`
	Path string `json:"path,omitempty"` // Default is be /promote in controller logic
}

// MasterServiceSpec defines the desired state of MasterService
// +k8s:openapi-gen=true
type MasterServiceSpec struct {
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html
	ServiceSpec corev1.ServiceSpec    `json:"serviceSpec"`
	Callback    MasterServiceCallback `json:"callback,omitempty"`
}

// MasterServiceStatus defines the observed state of MasterService
// +k8s:openapi-gen=true
type MasterServiceStatus struct {
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MasterService is the Schema for the masterservices API
// +k8s:openapi-gen=true
type MasterService struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MasterServiceSpec   `json:"spec"`
	Status MasterServiceStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MasterServiceList contains a list of MasterService
type MasterServiceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MasterService `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MasterService{}, &MasterServiceList{})
}
