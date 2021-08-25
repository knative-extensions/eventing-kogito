/*
Copyright 2019 The Knative Authors.

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

package v1alpha1

import (
	kogitoapi "github.com/kiegroup/kogito-operator/api/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"knative.dev/pkg/apis"
	"knative.dev/pkg/apis/duck"
	duckv1 "knative.dev/pkg/apis/duck/v1"
	"knative.dev/pkg/kmeta"
	"knative.dev/pkg/webhook/resourcesemantics"
)

// +genclient
// +genreconciler
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type KogitoSource struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec holds the desired state of the KogitoSource (from the client).
	Spec KogitoSourceSpec `json:"spec"`

	// Status communicates the observed state of the KogitoSource (from the controller).
	// +optional
	Status KogitoSourceStatus `json:"status,omitempty"`
}

// GetGroupVersionKind returns the GroupVersionKind.
func (*KogitoSource) GetGroupVersionKind() schema.GroupVersionKind {
	return SchemeGroupVersion.WithKind("KogitoSource")
}

var (
	// Check that KogitoSource can be validated and defaulted.
	_ apis.Validatable = (*KogitoSource)(nil)
	_ apis.Defaultable = (*KogitoSource)(nil)
	// Check that we can create OwnerReferences to a KogitoSource.
	_ kmeta.OwnerRefable = (*KogitoSource)(nil)
	// Check that KogitoSource is a runtime.Object.
	_ runtime.Object = (*KogitoSource)(nil)
	// Check that KogitoSource satisfies resourcesemantics.GenericCRD.
	_ resourcesemantics.GenericCRD = (*KogitoSource)(nil)
	// Check that KogitoSource implements the Conditions duck type.
	_ = duck.VerifyType(&KogitoSource{}, &duckv1.Conditions{})
	// Check that the type conforms to the duck Knative Resource shape.
	_ duckv1.KRShaped = (*KogitoSource)(nil)
)

// KogitoSourceSpec holds the desired state of the KogitoSource (from the client).
type KogitoSourceSpec struct {
	// inherits duck/v1 SourceSpec, which currently provides:
	// * Sink - a reference to an object that will resolve to a domain name or
	//   a URI directly to use as the sink.
	// * CloudEventOverrides - defines overrides to control the output format
	//   and modifications of the event sent to the sink.
	duckv1.SourceSpec `json:",inline"`

	// ServiceAccountName holds the name of the Kubernetes service account
	// as which the underlying K8s resources should be run. If unspecified
	// this will default to the "default" service account for the namespace
	// in which the KogitoSource exists.
	// +optional
	ServiceAccountName string `json:"serviceAccountName,omitempty"`

	// inherits kogitoapi.KogitoRuntimeSpec, which provides the interface for users
	// to declare the KogitoRuntime service to be deployed as an event source
	kogitoapi.KogitoRuntimeSpec `json:",inline"`
}

const (
	// KogitoSourceConditionReady is set when the revision is starting to materialize
	// runtime resources, and becomes true when those resources are ready.
	KogitoSourceConditionReady = apis.ConditionReady
)

// KogitoSourceStatus communicates the observed state of the KogitoSource (from the controller).
type KogitoSourceStatus struct {
	// inherits duck/v1 SourceStatus, which currently provides:
	// * ObservedGeneration - the 'Generation' of the Service that was last
	//   processed by the controller.
	// * Conditions - the latest available observations of a resource's current
	//   state.
	// * SinkURI - the current active sink URI that has been configured for the
	//   Source.
	duckv1.SourceStatus `json:",inline"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KogitoSourceList is a list of KogitoSource resources
type KogitoSourceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []KogitoSource `json:"items"`
}

// GetStatus retrieves the status of the resource. Implements the KRShaped interface.
func (ss *KogitoSource) GetStatus() *duckv1.Status {
	return &ss.Status.Status
}
