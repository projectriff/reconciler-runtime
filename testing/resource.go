/*
Copyright 2019 the original author or authors.

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

package testing

import (
	"fmt"
	"time"

	"github.com/projectriff/reconciler-runtime/apis"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

var _ webhook.Defaulter = &TestResource{}

// +kubebuilder:object:root=true
// +genclient

type TestResource struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TestResourceSpec   `json:"spec"`
	Status TestResourceStatus `json:"status"`
}

func (r *TestResource) Default() {
	if r.Spec.Fields == nil {
		r.Spec.Fields = map[string]string{}
	}
	r.Spec.Fields["Defaulter"] = "ran"
}

// +kubebuilder:object:generate=true
type TestResourceSpec struct {
	Fields map[string]string `json:"fields,omitempty"`
}

// +kubebuilder:object:generate=true
type TestResourceStatus struct {
	apis.Status `json:",inline"`
	Fields      map[string]string `json:"fields,omitempty"`
}

func (rs *TestResourceStatus) InitializeConditions() {
	condSet := apis.NewLivingConditionSet()
	condSet.Manage(rs).InitializeConditions()
}

func (rs *TestResourceStatus) MarkReady() {
	rs.SetConditions(apis.Conditions{
		{
			Type:               apis.ConditionReady,
			Status:             corev1.ConditionTrue,
			LastTransitionTime: apis.VolatileTime{Inner: metav1.NewTime(time.Now())},
		},
	})
}

func (rs *TestResourceStatus) MarkNotReady(reason, message string, messageA ...interface{}) {
	rs.SetConditions(apis.Conditions{
		{
			Type:               apis.ConditionReady,
			Status:             corev1.ConditionFalse,
			Reason:             reason,
			Message:            fmt.Sprintf(message, messageA...),
			LastTransitionTime: apis.VolatileTime{Inner: metav1.NewTime(time.Now())},
		},
	})
	condSet := apis.NewLivingConditionSet()
	condSet.Manage(rs).MarkFalse(apis.ConditionReady, reason, message, messageA...)
}

// +kubebuilder:object:root=true

type TestResourceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []TestResource `json:"items"`
}

var (
	// GroupVersion is group version used to register these objects
	GroupVersion = schema.GroupVersion{Group: "testing.reconciler.runtime", Version: "v1"}

	// SchemeBuilder is used to add go types to the GroupVersionKind scheme
	SchemeBuilder = &scheme.Builder{GroupVersion: GroupVersion}

	// AddToScheme adds the types in this group-version to the given scheme.
	AddToScheme = SchemeBuilder.AddToScheme
)

// compatibility with k8s.io/code-generator
var SchemeGroupVersion = GroupVersion

func init() {
	SchemeBuilder.Register(&TestResource{}, &TestResourceList{})
}
