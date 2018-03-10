/*
Copyright 2018 Jack Lin (jacklin@lsalab.cs.nthu.edu.tw)

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

// Package v1alpha1 for a sample crd
package v1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Student struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`

	Spec   StudentSpec   `json:"spec"`
    Status StudentStatus `json:"status"`
}

type StudentSpec struct {
	TaskName string `json:"taskName"`
    Threads  *int32 `json:"threads"`
}

type StudentLifeState string

const (
        Health  StudentLifeState = "Health!"
        Sick    StudentLifeState = "Sick!"
        Dead    StudentLifeState = "Dead!"

)

type StudentStatus struct {
    AvailableThreads int32              `json:"availableThreads"`
    LifeState        StudentLifeState   `json:"lifeState"`
    Message          string             `json:"message"`
    LastLiveTime     metav1.Time        `json:"lastLiveTime"`
}


// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type StudentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []Student `json:"items"`
}
