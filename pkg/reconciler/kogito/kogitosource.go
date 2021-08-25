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

package kogito

import (
	"context"

	// k8s.io imports
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	// knative.dev/pkg imports
	"knative.dev/pkg/logging"
	pkgreconciler "knative.dev/pkg/reconciler"
	"knative.dev/pkg/resolver"

	// knative.dev/eventing imports
	sourcesv1 "knative.dev/eventing/pkg/apis/sources/v1"
	reconcilersource "knative.dev/eventing/pkg/reconciler/source"

	"knative.dev/eventing-kogito/pkg/apis/kogito/v1alpha1"
	reconcilerkogitosource "knative.dev/eventing-kogito/pkg/client/injection/reconciler/kogito/v1alpha1/kogitosource"
	"knative.dev/eventing-kogito/pkg/reconciler"
	"knative.dev/eventing-kogito/pkg/reconciler/kogito/resources"
)

// Reconciler reconciles a KogitoSource object
type Reconciler struct {
	krr            *reconciler.KogitoRuntimeReconciler
	sinkResolver   *resolver.URIResolver
	configAccessor reconcilersource.ConfigAccessor
}

// Check that our Reconciler implements Interface
var _ reconcilerkogitosource.Interface = (*Reconciler)(nil)

// ReconcileKind implements Interface.ReconcileKind.
func (r *Reconciler) ReconcileKind(ctx context.Context, src *v1alpha1.KogitoSource) pkgreconciler.Event {

	ctx = sourcesv1.WithURIResolver(ctx, r.sinkResolver)

	ra, sb, event := r.krr.ReconcileKogitoRuntime(ctx, src, makeSinkBinding(src),
		resources.MakeReceiveAdapter(&resources.ReceiveAdapterArgs{
			EventSource:    src.Namespace + "/" + src.Name,
			Source:         src,
			Labels:         resources.Labels(src.Name),
			AdditionalEnvs: r.configAccessor.ToEnvVars(), // Grab config envs for tracing/logging/metrics
		}),
	)
	if ra != nil {
		src.Status.PropagateDeploymentAvailability(ra)
	}
	if sb != nil {
		if c := sb.Status.GetCondition(sourcesv1.SinkBindingConditionSinkProvided); c.IsTrue() {
			src.Status.MarkSink(sb.Status.SinkURI)
		} else if c.IsFalse() {
			src.Status.MarkNoSink(c.GetReason(), "%s", c.GetMessage())
		}
	}
	if event != nil {
		logging.FromContext(ctx).Infof("returning because event from ReconcileKogitoRuntime")
		return event
	}

	return nil
}

func makeSinkBinding(src *v1alpha1.KogitoSource) *sourcesv1.SinkBinding {
	return &sourcesv1.SinkBinding{
		ObjectMeta: metav1.ObjectMeta{
			// this is necessary to track the change of sink reference.
			Name:      src.GetName(),
			Namespace: src.GetNamespace(),
		},
		Spec: sourcesv1.SinkBindingSpec{
			SourceSpec: src.Spec.SourceSpec,
		},
	}
}
