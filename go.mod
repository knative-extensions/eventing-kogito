module knative.dev/eventing-kogito

go 1.14

require (
	github.com/google/go-cmp v0.5.6
	go.uber.org/zap v1.19.1
	k8s.io/api v0.21.4
	k8s.io/apimachinery v0.21.4
	k8s.io/client-go v0.21.4
	knative.dev/eventing v0.28.1-0.20211217092418-fede720191d3
	knative.dev/hack v0.0.0-20211216134818-6fc030496333
	knative.dev/pkg v0.0.0-20211216142117-79271798f696
)

replace github.com/prometheus/client_golang => github.com/prometheus/client_golang v0.9.2
