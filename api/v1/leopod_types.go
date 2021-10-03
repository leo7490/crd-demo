/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// LeoPodSpec defines the desired state of LeoPod
type LeoPodSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of LeoPod. Edit leopod_types.go to remove/update
	// Foo      string `json:"foo,omitempty"`
	Image    string `json:"image"`
	Replicas int32  `json:"replicas"`
}

// LeoPodStatus defines the observed state of LeoPod
type LeoPodStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	ReadyReplicas int32 `json:"readyReplicas,omitempty" protobuf:"varint,7,opt,name=readyReplicas"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// LeoPod is the Schema for the leopods API
type LeoPod struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   LeoPodSpec   `json:"spec,omitempty"`
	Status LeoPodStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// LeoPodList contains a list of LeoPod
type LeoPodList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []LeoPod `json:"items"`
}

func init() {
	SchemeBuilder.Register(&LeoPod{}, &LeoPodList{})
}
