module knative.dev/eventing-kogito

go 1.14

require (
	github.com/google/go-cmp v0.5.6
	go.uber.org/zap v1.21.0
	k8s.io/api v0.24.4
	k8s.io/apimachinery v0.24.4
	k8s.io/client-go v0.24.4
	knative.dev/eventing v0.34.1-0.20220919143309-645fd887ce0c
	knative.dev/hack v0.0.0-20220914183605-d1317b08c0c3
	knative.dev/pkg v0.0.0-20220920224309-df29e2a20aba
)

replace github.com/prometheus/client_golang => github.com/prometheus/client_golang v0.9.2
