module knative.dev/eventing-kogito

go 1.14

require (
	github.com/google/go-cmp v0.5.6
	go.uber.org/zap v1.21.0
	k8s.io/api v0.23.9
	k8s.io/apimachinery v0.23.9
	k8s.io/client-go v0.23.9
	knative.dev/eventing v0.33.1-0.20220822135655-5ee615866b28
	knative.dev/hack v0.0.0-20220815132133-e9a8475f4329
	knative.dev/pkg v0.0.0-20220818004048-4a03844c0b15
)

replace github.com/prometheus/client_golang => github.com/prometheus/client_golang v0.9.2
