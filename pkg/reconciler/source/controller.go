package source

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"knative.dev/eventing-kogito/pkg/apis/kogito/v1alpha1"
	"knative.dev/eventing-kogito/pkg/client/clientset/versioned/scheme"
	kogitosourceinformer "knative.dev/eventing-kogito/pkg/client/injection/informers/kogito/v1alpha1/kogitosource"
	v1 "knative.dev/eventing/pkg/apis/sources/v1"
	"knative.dev/pkg/apis/duck"
	"knative.dev/pkg/client/injection/ducks/duck/v1/podspecable"
	"knative.dev/pkg/configmap"
	"knative.dev/pkg/controller"
	"knative.dev/pkg/injection/clients/dynamicclient"
	"knative.dev/pkg/logging"
	"knative.dev/pkg/resolver"
	"knative.dev/pkg/tracker"
	"knative.dev/pkg/webhook/psbinding"
)

const controllerAgentName = "kogitosourcebinding-controller"

// NewController for the KogitoSource Bindable interface
// See: https://github.com/knative/pkg/tree/main/webhook/psbinding for reference
func NewController(ctx context.Context, cmw configmap.Watcher) *controller.Impl {
	logger := logging.FromContext(ctx)

	kogitoInformer := kogitosourceinformer.Get(ctx)
	dc := dynamicclient.Get(ctx)
	psInformerFactory := podspecable.Get(ctx)

	c := &psbinding.BaseReconciler{
		GVR: v1alpha1.SchemeGroupVersion.WithResource("kogitosources"),
		Get: func(namespace string, name string) (psbinding.Bindable, error) {
			return kogitoInformer.Lister().KogitoSources(namespace).Get(name)
		},
		DynamicClient: dc,
		Recorder: record.NewBroadcaster().NewRecorder(
			scheme.Scheme, corev1.EventSource{Component: controllerAgentName}),
	}
	logger = logger.Named("KogitoSourceBindings")
	impl := controller.NewContext(ctx, c, controller.ControllerOptions{WorkQueueName: "KogitoSourceBinding", Logger: logger})

	sbResolver := resolver.NewURIResolverFromTracker(ctx, impl.Tracker)
	c.SubResourcesReconciler = &KogitoSourceSubResourcesReconciler{
		tracker: impl.Tracker,
		res:     sbResolver,
	}

	logger.Info("Setting up event handlers")

	kogitoInformer.Informer().AddEventHandler(controller.HandleAll(impl.Enqueue))

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

func NewWebhook(ctx context.Context, cmw configmap.Watcher) *controller.Impl {
	return psbinding.NewAdmissionController(ctx,
		// Name of the resource webhook.
		"kogitosources.webhook.bindings.kogito.kie.org",

		// The path on which to serve the webhook.
		"/kogitosources",

		// How to get all the Bindables for configuring the mutating webhook.
		ListAll,

		// How to setup the context prior to invoking Do/Undo.
		func(ctx context.Context, b psbinding.Bindable) (context.Context, error) {
			return ctx, nil
		},
	)
}

// ListAll ...
func ListAll(ctx context.Context, handler cache.ResourceEventHandler) psbinding.ListAll {
	kogitoInformer := kogitosourceinformer.Get(ctx)

	// Whenever a KogitoSource changes our webhook programming might change.
	kogitoInformer.Informer().AddEventHandler(handler)

	return func() ([]psbinding.Bindable, error) {
		l, err := kogitoInformer.Lister().List(labels.Everything())
		if err != nil {
			return nil, err
		}
		bl := make([]psbinding.Bindable, 0, len(l))
		for _, elt := range l {
			bl = append(bl, elt)
		}
		return bl, nil
	}

}
