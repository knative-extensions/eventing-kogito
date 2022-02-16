module knative.dev/eventing-kogito

go 1.14

require (
	github.com/google/go-cmp v0.5.6
	go.uber.org/zap v1.19.1
	k8s.io/api v0.22.5
	k8s.io/apimachinery v0.22.5
	k8s.io/client-go v0.22.5
	knative.dev/eventing v0.29.1-0.20220214155047-0326f92dd5b7
	knative.dev/hack v0.0.0-20220215185059-b9cb1983b600
	knative.dev/pkg v0.0.0-20220215153400-3c00bb0157b9
)

replace github.com/prometheus/client_golang => github.com/prometheus/client_golang v0.9.2
