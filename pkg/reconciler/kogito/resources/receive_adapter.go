/*
Copyright 2019 The Knative Authors

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

package resources

import (
	"github.com/kiegroup/kogito-operator/api/v1beta1"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knative.dev/pkg/kmeta"

	"knative.dev/eventing-kogito/pkg/apis/kogito/v1alpha1"
)

// ReceiveAdapterArgs are the arguments needed to create a KogitoSource Receive Adapter.
// Every field is required.
type ReceiveAdapterArgs struct {
	Labels         map[string]string
	Source         *v1alpha1.KogitoSource
	EventSource    string
	AdditionalEnvs []corev1.EnvVar
}

// MakeReceiveAdapter generates (but does not insert into K8s) the Receive Adapter KogitoRuntime for
// KogitoSource sources.
func MakeReceiveAdapter(args *ReceiveAdapterArgs) *v1beta1.KogitoRuntime {
	replicas := int32(1)
	kogitoRuntime := &v1beta1.KogitoRuntime{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: args.Source.Namespace,
			Name:      kmeta.ChildName("ks-", args.Source.Name),
			Labels:    args.Labels,
			OwnerReferences: []metav1.OwnerReference{
				*kmeta.NewControllerRef(args.Source),
			},
		},
		Spec: args.Source.Spec.KogitoRuntimeSpec,
	}
	if kogitoRuntime.Spec.Replicas == nil {
		kogitoRuntime.Spec.Replicas = &replicas
	}
	kogitoRuntime.Spec.DeploymentLabels = args.Labels
	kogitoRuntime.Spec.ServiceLabels = args.Labels
	kogitoRuntime.Spec.Env = append(makeEnv(args.EventSource), args.AdditionalEnvs...)
	return kogitoRuntime
}

func makeEnv(eventSource string) []corev1.EnvVar {
	return []corev1.EnvVar{{
		Name:  "EVENT_SOURCE",
		Value: eventSource,
	}, {
		Name:  "METRICS_DOMAIN",
		Value: "knative.dev/eventing",
	}}
}
