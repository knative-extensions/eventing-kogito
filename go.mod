module knative.dev/eventing-kogito

go 1.14

require (
	github.com/google/go-cmp v0.5.6
	go.uber.org/zap v1.19.1
	k8s.io/api v0.22.5
	k8s.io/apimachinery v0.22.5
	k8s.io/client-go v0.22.5
	knative.dev/eventing v0.29.1-0.20220217054812-7a48f4269b6f
	knative.dev/hack v0.0.0-20220218190734-a8ef7b67feec
	knative.dev/pkg v0.0.0-20220217155112-d48172451966
)

replace github.com/prometheus/client_golang => github.com/prometheus/client_golang v0.9.2
