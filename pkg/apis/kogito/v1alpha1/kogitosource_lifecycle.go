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
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "knative.dev/eventing/pkg/apis/sources/v1"
	"knative.dev/pkg/apis"
	duckapi "knative.dev/pkg/apis/duck"
	duckv1 "knative.dev/pkg/apis/duck/v1"
	"knative.dev/pkg/tracker"
)

const (
	// KogitoConditionReady has status True when the KogitoSource is ready to send events.
	KogitoConditionReady = apis.ConditionReady

	// KogitoConditionSinkProvided has status True when the KogitoSource has been configured with a sink target.
	KogitoConditionSinkProvided apis.ConditionType = "SinkProvided"

	// KogitoConditionDeployed has status True when the KogitoSource has had it's deployment created.
	KogitoConditionDeployed apis.ConditionType = "Deployed"
)

var KogitoCondSet = apis.NewLivingConditionSet(
	KogitoConditionSinkProvided,
	KogitoConditionDeployed,
)

// GetConditionSet returns KogitoSource ConditionSet.
func (*KogitoSource) GetConditionSet() apis.ConditionSet {
	return KogitoCondSet
}

func (ks *KogitoSource) GetSubject() tracker.Reference {
	return ks.Spec.Subject
}

func (ks *KogitoSource) GetBindingStatus() duckapi.BindableStatus {
	return &ks.Status
}

// Do reuse the logic from SinkBinding to inject the sink URLs to the target pod
func (ks *KogitoSource) Do(ctx context.Context, pod *duckv1.WithPod) {
	// overload SinkBinding's Do
	ks.sinkBinding().Do(ctx, pod)
	// other actions with the bindable object ...
}

// Undo reuse the logic from SinkBinding to remove the injected environment variables from the target pod
func (ks *KogitoSource) Undo(ctx context.Context, pod *duckv1.WithPod) {
	// overload SinkBinding's Undo
	ks.sinkBinding().Undo(ctx, pod)
	// other actions with the bindable object ...
}

func (ks *KogitoSource) sinkBinding() *v1.SinkBinding {
	return &v1.SinkBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ks.Name,
			Namespace: ks.Namespace,
		},
		Spec: v1.SinkBindingSpec{
			SourceSpec: ks.Spec.SourceSpec,
		},
	}
}

// SetObservedGeneration implements psbinding.BindableStatus
func (sbs *KogitoSourceStatus) SetObservedGeneration(gen int64) {
	sbs.ObservedGeneration = gen
}

// MarkBindingUnavailable marks the KogitoSource's Ready condition to False with
// the provided reason and message.
func (sbs *KogitoSourceStatus) MarkBindingUnavailable(reason, message string) {
	KogitoCondSet.Manage(sbs).MarkFalse(KogitoConditionReady, reason, message)
}

// MarkBindingAvailable marks the KogitoSource's Ready condition to True.
func (sbs *KogitoSourceStatus) MarkBindingAvailable() {
	KogitoCondSet.Manage(sbs).MarkTrue(KogitoConditionReady)
}

// GetCondition returns the condition currently associated with the given type, or nil.
func (s *KogitoSourceStatus) GetCondition(t apis.ConditionType) *apis.Condition {
	return KogitoCondSet.Manage(s).GetCondition(t)
}

// InitializeConditions sets relevant unset conditions to Unknown state.
func (s *KogitoSourceStatus) InitializeConditions() {
	KogitoCondSet.Manage(s).InitializeConditions()
}

// MarkSink sets the condition that the source has a sink configured.
func (s *KogitoSourceStatus) MarkSink(uri *apis.URL) {
	s.SinkURI = uri
	if len(uri.String()) > 0 {
		KogitoCondSet.Manage(s).MarkTrue(KogitoConditionSinkProvided)
	} else {
		KogitoCondSet.Manage(s).MarkUnknown(KogitoConditionSinkProvided, "SinkEmpty", "Sink has resolved to empty.")
	}
}

// MarkNoSink sets the condition that the source does not have a sink configured.
func (s *KogitoSourceStatus) MarkNoSink(reason, messageFormat string, messageA ...interface{}) {
	s.SinkURI = nil
	KogitoCondSet.Manage(s).MarkFalse(KogitoConditionSinkProvided, reason, messageFormat, messageA...)
}

// IsReady returns true if the resource is ready overall.
func (s *KogitoSourceStatus) IsReady() bool {
	return KogitoCondSet.Manage(s).IsHappy()
}
