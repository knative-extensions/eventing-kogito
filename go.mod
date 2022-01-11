module knative.dev/eventing-kogito

go 1.14

require (
	github.com/google/go-cmp v0.5.6
	go.uber.org/zap v1.19.1
	k8s.io/api v0.22.5
	k8s.io/apimachinery v0.22.5
	k8s.io/client-go v0.22.5
	knative.dev/eventing v0.28.1-0.20220111105413-b5603c0ad63d
	knative.dev/hack v0.0.0-20220110200259-f08cb0dcdee7
	knative.dev/pkg v0.0.0-20220105211333-96f18522d78d
)

replace github.com/prometheus/client_golang => github.com/prometheus/client_golang v0.9.2
