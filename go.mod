module knative.dev/eventing-kogito

go 1.14

require (
	github.com/google/go-cmp v0.5.6
	go.uber.org/zap v1.21.0
	k8s.io/api v0.23.8
	k8s.io/apimachinery v0.23.8
	k8s.io/client-go v0.23.8
	knative.dev/eventing v0.32.1-0.20220630063730-159f4ec31429
	knative.dev/hack v0.0.0-20220629134730-e7d63651ce8f
	knative.dev/pkg v0.0.0-20220630112730-85965e1e8eb1
)

replace github.com/prometheus/client_golang => github.com/prometheus/client_golang v0.9.2
