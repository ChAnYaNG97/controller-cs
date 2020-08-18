/*
Copyright 2020 yangchen.

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

type PhaseType string

const (
	Phase_None   PhaseType = ""
	Phase_Create PhaseType = "Creating"
	Phase_Run    PhaseType = "Running"
	Phase_Delete PhaseType = "Deleting"
	Phase_Error  PhaseType = "Error"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// AliyunCKSpec defines the desired state of AliyunCK
type AliyunCKSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Region
	Region			string `json:"region"`

	// Machine Type
	InstanceType	string `json:"instance_type"`
}

// AliyunCKStatus defines the observed state of AliyunCK
type AliyunCKStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file


	// Status
	Phase		PhaseType `json:"status"`

	// Cluster Id
	ClusterId	string `json:"cluster_id"`
	VPCId		string `json:"vpc_id"`

}

// +kubebuilder:object:root=true

// AliyunCK is the Schema for the aliyuncks API
type AliyunCK struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AliyunCKSpec   `json:"spec,omitempty"`
	Status AliyunCKStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// AliyunCKList contains a list of AliyunCK
type AliyunCKList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AliyunCK `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AliyunCK{}, &AliyunCKList{})
}
