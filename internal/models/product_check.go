package models

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// ProductCheck defines the structure of the ProductCheck CRD.

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
type ProductCheck struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              ProductCheckSpec   `json:"spec,omitempty"`
	Status            ProductCheckStatus `json:"status,omitempty"`
}

// ProductCheckSpec defines the spec of the ProductCheck CRD.
type ProductCheckSpec struct {
	ProductName string `json:"productName"`
	Version     string `json:"version"`
}

// ProductCheckStatus defines the status of the ProductCheck CRD.
type ProductCheckStatus struct {
	EndOfLifeDate string `json:"endOfLifeDate"`
	Status        string `json:"status"`
}

// ProductCheckList defines a list of ProductCheck objects.
// +kubebuilder:object:root=true
type ProductCheckList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ProductCheck `json:"items"`
}

// DeepCopyObject implements client.ObjectList.
func (p *ProductCheckList) DeepCopyObject() runtime.Object {
	panic("unimplemented")
}

// GetContinue implements client.ObjectList.
// Subtle: this method shadows the method (ListMeta).GetContinue of ProductCheckList.ListMeta.
func (p *ProductCheckList) GetContinue() string {
	panic("unimplemented")
}

// GetObjectKind implements client.ObjectList.
// Subtle: this method shadows the method (TypeMeta).GetObjectKind of ProductCheckList.TypeMeta.
func (p *ProductCheckList) GetObjectKind() schema.ObjectKind {
	panic("unimplemented")
}

// GetRemainingItemCount implements client.ObjectList.
// Subtle: this method shadows the method (ListMeta).GetRemainingItemCount of ProductCheckList.ListMeta.
func (p *ProductCheckList) GetRemainingItemCount() *int64 {
	panic("unimplemented")
}

// GetResourceVersion implements client.ObjectList.
// Subtle: this method shadows the method (ListMeta).GetResourceVersion of ProductCheckList.ListMeta.
func (p *ProductCheckList) GetResourceVersion() string {
	panic("unimplemented")
}

// GetSelfLink implements client.ObjectList.
// Subtle: this method shadows the method (ListMeta).GetSelfLink of ProductCheckList.ListMeta.
func (p *ProductCheckList) GetSelfLink() string {
	panic("unimplemented")
}

// SetContinue implements client.ObjectList.
// Subtle: this method shadows the method (ListMeta).SetContinue of ProductCheckList.ListMeta.
func (p *ProductCheckList) SetContinue(c string) {
	panic("unimplemented")
}

// SetRemainingItemCount implements client.ObjectList.
// Subtle: this method shadows the method (ListMeta).SetRemainingItemCount of ProductCheckList.ListMeta.
func (p *ProductCheckList) SetRemainingItemCount(c *int64) {
	panic("unimplemented")
}

// SetResourceVersion implements client.ObjectList.
// Subtle: this method shadows the method (ListMeta).SetResourceVersion of ProductCheckList.ListMeta.
func (p *ProductCheckList) SetResourceVersion(version string) {
	panic("unimplemented")
}

// SetSelfLink implements client.ObjectList.
// Subtle: this method shadows the method (ListMeta).SetSelfLink of ProductCheckList.ListMeta.
func (p *ProductCheckList) SetSelfLink(selfLink string) {
	panic("unimplemented")
}
