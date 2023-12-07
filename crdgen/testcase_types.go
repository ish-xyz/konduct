package crdgen

import (
	"github.com/ish-xyz/ykubetest/pkg/loader"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:object:root=true
type TestCase struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              loader.TestCase `json:"spec,omitempty"`
	Status            TestCaseStatus  `json:"status,omitempty"`
}

type TestCaseStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// TestCaseList contains a list of TestCase
type TestCaseList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TestCase `json:"items"`
}

type Definition struct {
	// +kubebuilder:validation:Required
	Group string `json:"apiGroup"`
	// +kubebuilder:validation:Required
	Resource string `json:"resource"`
	// +kubebuilder:validation:Required
	Operation string `json:"operation"`
}
