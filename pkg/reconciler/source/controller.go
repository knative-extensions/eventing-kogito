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

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"knative.dev/eventing-kogito/pkg/apis/kogito/v1alpha1"
	"knative.dev/eventing-kogito/pkg/client/clientset/versioned/scheme"
	kogitosourceinformer "knative.dev/eventing-kogito/pkg/client/injection/informers/kogito/v1alpha1/kogitosource"
	v1 "knative.dev/eventing/pkg/apis/sources/v1"
	"knative.dev/pkg/apis/duck"
	"knative.dev/pkg/client/injection/ducks/duck/v1/podspecable"
	"knative.dev/pkg/client/injection/kube/informers/core/v1/namespace"
	"knative.dev/pkg/configmap"
	"knative.dev/pkg/controller"
	"knative.dev/pkg/injection/clients/dynamicclient"
	"knative.dev/pkg/logging"
	"knative.dev/pkg/reconciler"
	"knative.dev/pkg/resolver"
	"knative.dev/pkg/tracker"
	"knative.dev/pkg/webhook/psbinding"
)

const controllerAgentName = "kogitosource-controller"

// NewController for the KogitoSource Bindable interface
// See: https://github.com/knative/pkg/tree/main/webhook/psbinding for reference
func NewController(ctx context.Context, cmw configmap.Watcher) *controller.Impl {
	logger := logging.FromContext(ctx)

	kogitoInformer := kogitosourceinformer.Get(ctx)
	dc := dynamicclient.Get(ctx)
	psInformerFactory := podspecable.Get(ctx)
	namespaceInformer := namespace.Get(ctx)

	c := &psbinding.BaseReconciler{
		LeaderAwareFuncs: reconciler.LeaderAwareFuncs{
			PromoteFunc: func(bkt reconciler.Bucket, enq func(reconciler.Bucket, types.NamespacedName)) error {
				all, err := kogitoInformer.Lister().List(labels.Everything())
				if err != nil {
					return err
				}
				for _, elt := range all {
					enq(bkt, types.NamespacedName{
						Namespace: elt.GetNamespace(),
						Name:      elt.GetName(),
					})
				}
				return nil
			},
		},
		GVR: v1alpha1.SchemeGroupVersion.WithResource("kogitosources"),
		Get: func(namespace string, name string) (psbinding.Bindable, error) {
			return kogitoInformer.Lister().KogitoSources(namespace).Get(name)
		},
		DynamicClient:   dc,
		NamespaceLister: namespaceInformer.Lister(),
		Recorder: record.NewBroadcaster().NewRecorder(
			scheme.Scheme, corev1.EventSource{Component: controllerAgentName}),
	}
	logger = logger.Named("KogitoSource")
	impl := controller.NewContext(ctx, c, controller.ControllerOptions{WorkQueueName: "KogitoSource", Logger: logger})

	sbResolver := resolver.NewURIResolverFromTracker(ctx, impl.Tracker)
	c.SubResourcesReconciler = &KogitoSourceSubResourcesReconciler{
		tracker: impl.Tracker,
		res:     sbResolver,
	}

	logger.Info("Setting up event handlers")

	kogitoInformer.Informer().AddEventHandler(controller.HandleAll(impl.Enqueue))
	namespaceInformer.Informer().AddEventHandler(controller.HandleAll(impl.Enqueue))

	c.Tracker = tracker.New(impl.EnqueueKey, controller.GetTrackerLease(ctx))
	c.Factory = &duck.CachedInformerFactory{
		Delegate: &duck.EnqueueInformerFactory{
			Delegate:     psInformerFactory,
			EventHandler: controller.HandleAll(c.Tracker.OnChanged),
		},
	}

	// If our `Do` / `Undo` methods need additional context, then we can
	// setup a callback to infuse the `context.Context` here:
	//    c.WithContext = ...
	// Note that this can also set up additional informer watch events to
	// trigger reconciliation when the infused context changes.
	c.WithContext = func(ctx context.Context, b psbinding.Bindable) (context.Context, error) {
		return v1.WithURIResolver(ctx, sbResolver), nil
	}

	return impl
}
