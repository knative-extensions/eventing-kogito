module knative.dev/eventing-kogito

go 1.14

require (
	github.com/google/go-cmp v0.5.6
	go.uber.org/zap v1.21.0
	k8s.io/api v0.23.8
	k8s.io/apimachinery v0.23.8
	k8s.io/client-go v0.23.8
	knative.dev/eventing v0.33.1-0.20220722122720-c8435ed74ba8
	knative.dev/hack v0.0.0-20220725145124-782bbaabb8a1
	knative.dev/pkg v0.0.0-20220722175921-6c9c1c6098d5
)

replace github.com/prometheus/client_golang => github.com/prometheus/client_golang v0.9.2
