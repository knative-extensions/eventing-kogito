module knative.dev/eventing-kogito

go 1.14

require (
	github.com/google/go-cmp v0.5.6
	go.uber.org/zap v1.19.1
	k8s.io/api v0.22.5
	k8s.io/apimachinery v0.22.5
	k8s.io/client-go v0.22.5
	knative.dev/eventing v0.28.1-0.20220113075012-61d97d366c26
	knative.dev/hack v0.0.0-20220111151514-59b0cf17578e
	knative.dev/pkg v0.0.0-20220113045912-c0e1594c2fb1
)

replace github.com/prometheus/client_golang => github.com/prometheus/client_golang v0.9.2
