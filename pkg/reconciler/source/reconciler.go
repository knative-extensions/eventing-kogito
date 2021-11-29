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

package source

import (
	"context"
	"errors"

	"knative.dev/eventing-kogito/pkg/apis/kogito/v1alpha1"
	"knative.dev/pkg/logging"
	"knative.dev/pkg/resolver"
	"knative.dev/pkg/tracker"
	"knative.dev/pkg/webhook/psbinding"
)

type KogitoSourceSubResourcesReconciler struct {
	tracker tracker.Interface
	res     *resolver.URIResolver
}

// Reconcile we try to discover the Sink URI to define it in our Status
func (kbsrr *KogitoSourceSubResourcesReconciler) Reconcile(ctx context.Context, b psbinding.Bindable) error {
	kogitoSource := b.(*v1alpha1.KogitoSource)
	if kbsrr.res == nil {
		err := errors.New("KogitoSource Resolver is nil")
		logging.FromContext(ctx).Errorf("%w", err)
		kogitoSource.Status.MarkBindingUnavailable("NoResolver", "No resolver associated with context for sink")
		return err
	}
	if kogitoSource.Spec.Sink.Ref != nil {
		err := kbsrr.tracker.TrackReference(tracker.Reference{
			APIVersion: kogitoSource.Spec.Sink.Ref.APIVersion,
			Kind:       kogitoSource.Spec.Sink.Ref.Kind,
			Namespace:  kogitoSource.Spec.Sink.Ref.Namespace,
			Name:       kogitoSource.Spec.Sink.Ref.Name,
		}, b)
		if err != nil {
			logging.FromContext(ctx).Errorf("Failed to track Sink reference: %w", err)
			kogitoSource.Status.MarkBindingUnavailable("TrackFailed", "Error when tried to track Sink reference ")
			return err
		}
	}
	uri, err := kbsrr.res.URIFromDestinationV1(ctx, kogitoSource.Spec.Sink, kogitoSource)
	if err != nil {
		logging.FromContext(ctx).Errorf("Failed to get URI from Destination: %w", err)
		kogitoSource.Status.MarkBindingUnavailable("NoURI", "URI could not be extracted from destination ")
		return err
	}
	kogitoSource.Status.MarkSink(uri)
	return nil
}

func (kbsrr *KogitoSourceSubResourcesReconciler) ReconcileDeletion(ctx context.Context, b psbinding.Bindable) error {
	// Logic to delete k8s resources related to our Bindable
	// we are not creating anything, so no reason to delete
	return nil
}
