module knative.dev/eventing-kogito

go 1.14

require (
	github.com/google/go-cmp v0.5.6
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/kiegroup/kogito-operator/apis v0.0.0
	github.com/kiegroup/kogito-operator/client v0.0.0-20210913124655-178e5d4b8327
	github.com/stretchr/testify v1.7.0
	go.uber.org/zap v1.19.1
	k8s.io/api v0.21.4
	k8s.io/apimachinery v0.21.4
	k8s.io/client-go v0.21.4
	knative.dev/eventing v0.27.1-0.20211117135535-e7452e66b37c
	knative.dev/hack v0.0.0-20211117134436-69a2295d54ce
	knative.dev/pkg v0.0.0-20211117215328-5708c4c44232
)

replace (
	github.com/kiegroup/kogito-operator/apis => github.com/kiegroup/kogito-operator/apis v0.0.0-20210913124655-178e5d4b8327
	github.com/prometheus/client_golang => github.com/prometheus/client_golang v0.9.2
)
