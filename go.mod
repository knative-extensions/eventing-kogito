module knative.dev/eventing-kogito

go 1.14

require (
	github.com/google/go-cmp v0.5.6
	go.uber.org/zap v1.19.1
	k8s.io/api v0.23.5
	k8s.io/apimachinery v0.23.5
	k8s.io/client-go v0.23.5
	knative.dev/eventing v0.31.1-0.20220427013821-2d4784764d75
	knative.dev/hack v0.0.0-20220427014036-5f473869d377
	knative.dev/pkg v0.0.0-20220427013826-1f681e126af6
)

replace github.com/prometheus/client_golang => github.com/prometheus/client_golang v0.9.2
