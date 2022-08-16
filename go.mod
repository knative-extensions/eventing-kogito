module knative.dev/eventing-kogito

go 1.14

require (
	github.com/google/go-cmp v0.5.6
	go.uber.org/zap v1.21.0
	k8s.io/api v0.23.9
	k8s.io/apimachinery v0.23.9
	k8s.io/client-go v0.23.9
	knative.dev/eventing v0.33.1-0.20220815185049-ab981d6bba92
	knative.dev/hack v0.0.0-20220815132133-e9a8475f4329
	knative.dev/pkg v0.0.0-20220815215248-d02dcd0b0391
)

replace github.com/prometheus/client_golang => github.com/prometheus/client_golang v0.9.2
