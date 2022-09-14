module knative.dev/eventing-kogito

go 1.14

require (
	github.com/google/go-cmp v0.5.6
	go.uber.org/zap v1.21.0
	k8s.io/api v0.24.4
	k8s.io/apimachinery v0.24.4
	k8s.io/client-go v0.24.4
	knative.dev/eventing v0.34.1-0.20220912074833-e46f4d509d74
	knative.dev/hack v0.0.0-20220913095247-7556452c2b54
	knative.dev/pkg v0.0.0-20220913150450-501fbd54de3d
)

replace github.com/prometheus/client_golang => github.com/prometheus/client_golang v0.9.2
