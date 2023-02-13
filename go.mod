module knative.dev/eventing-kogito

go 1.14

require (
	github.com/google/go-cmp v0.5.6
	go.uber.org/zap v1.21.0
	k8s.io/api v0.23.9
	k8s.io/apimachinery v0.23.9
	k8s.io/client-go v0.23.9
	knative.dev/eventing v0.34.7
	knative.dev/hack v0.0.0-20220823140917-8d1e4ccf9dc3
	knative.dev/pkg v0.0.0-20221123011941-9d7bd235ceed
)

replace github.com/prometheus/client_golang => github.com/prometheus/client_golang v0.9.2
