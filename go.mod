module knative.dev/eventing-kogito

go 1.14

require (
	github.com/google/go-cmp v0.5.6
	go.uber.org/zap v1.21.0
	k8s.io/api v0.24.4
	k8s.io/apimachinery v0.24.4
	k8s.io/client-go v0.24.4
	knative.dev/eventing v0.34.1-0.20221006135132-d1827fc34a20
	knative.dev/hack v0.0.0-20221006013832-52a87e1514cb
	knative.dev/pkg v0.0.0-20221006211431-fb2e4fbaa477
)

replace github.com/prometheus/client_golang => github.com/prometheus/client_golang v0.9.2
