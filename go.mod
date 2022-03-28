module knative.dev/eventing-kogito

go 1.14

require (
	github.com/google/go-cmp v0.5.6
	go.uber.org/zap v1.19.1
	k8s.io/api v0.23.5
	k8s.io/apimachinery v0.23.5
	k8s.io/client-go v0.23.5
	knative.dev/eventing v0.30.1-0.20220325083448-4be06cdde807
	knative.dev/hack v0.0.0-20220318020218-14f832e506f8
	knative.dev/pkg v0.0.0-20220325200448-1f7514acd0c2
)

replace github.com/prometheus/client_golang => github.com/prometheus/client_golang v0.9.2
