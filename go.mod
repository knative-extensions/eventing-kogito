module knative.dev/eventing-kogito

go 1.14

require (
	github.com/google/go-cmp v0.5.6
	go.uber.org/zap v1.21.0
	k8s.io/api v0.24.4
	k8s.io/apimachinery v0.24.4
	k8s.io/client-go v0.24.4
	knative.dev/eventing v0.34.1-0.20220926080258-70974588d7d4
	knative.dev/hack v0.0.0-20220923094413-9b7638704a22
	knative.dev/pkg v0.0.0-20220921024409-d1d5c849073b
)

replace github.com/prometheus/client_golang => github.com/prometheus/client_golang v0.9.2
