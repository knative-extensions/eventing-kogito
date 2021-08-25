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
	appsv1 "k8s.io/api/apps/v1"
	"knative.dev/eventing/pkg/apis/duck"
	"knative.dev/pkg/apis"
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

// GetCondition returns the condition currently associated with the given type, or nil.
func (s *KogitoSourceStatus) GetCondition(t apis.ConditionType) *apis.Condition {
	return KogitoCondSet.Manage(s).GetCondition(t)
}

// InitializeConditions sets relevant unset conditions to Unknown state.
func (s *KogitoSourceStatus) InitializeConditions() {
	KogitoCondSet.Manage(s).InitializeConditions()
}

// GetConditionSet returns KogitoSource ConditionSet.
func (*KogitoSource) GetConditionSet() apis.ConditionSet {
	return KogitoCondSet
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

// PropagateDeploymentAvailability uses the availability of the provided Deployment to determine if
// KogitoConditionDeployed should be marked as true or false.
func (s *KogitoSourceStatus) PropagateDeploymentAvailability(d *appsv1.Deployment) {
	if duck.DeploymentIsAvailable(&d.Status, false) {
		KogitoCondSet.Manage(s).MarkTrue(KogitoConditionDeployed)
	} else {
		// I don't know how to propagate the status well, so just give the name of the Deployment
		// for now.
		KogitoCondSet.Manage(s).MarkFalse(KogitoConditionDeployed, "KogitoRuntimeUnavailable", "The KogitoRuntime '%s' is unavailable.", d.Name)
	}
}

// IsReady returns true if the resource is ready overall.
func (s *KogitoSourceStatus) IsReady() bool {
	return KogitoCondSet.Manage(s).IsHappy()
}
