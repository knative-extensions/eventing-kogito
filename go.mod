module knative.dev/eventing-kogito

go 1.14

require (
	github.com/emicklei/go-restful v2.9.5+incompatible // indirect
	github.com/google/go-cmp v0.5.8
	go.uber.org/zap v1.21.0
	k8s.io/api v0.25.2
	k8s.io/apimachinery v0.25.2
	k8s.io/client-go v0.25.2
	knative.dev/eventing v0.35.0
	knative.dev/hack v0.0.0-20221010154335-3fdc50b9c24a
	knative.dev/pkg v0.0.0-20221011175852-714b7630a836
)

replace github.com/prometheus/client_golang => github.com/prometheus/client_golang v0.9.2
