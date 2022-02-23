module knative.dev/eventing-kogito

go 1.14

require (
	github.com/google/go-cmp v0.5.6
	go.uber.org/zap v1.19.1
	k8s.io/api v0.22.5
	k8s.io/apimachinery v0.22.5
	k8s.io/client-go v0.22.5
	knative.dev/eventing v0.29.1-0.20220221175003-2a69eec1e358
	knative.dev/hack v0.0.0-20220222192704-cf8cbc0e9165
	knative.dev/pkg v0.0.0-20220222211204-80c511aa340f
)

replace github.com/prometheus/client_golang => github.com/prometheus/client_golang v0.9.2
