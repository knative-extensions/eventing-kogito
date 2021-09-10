module knative.dev/eventing-kogito

go 1.14

require (
	github.com/google/go-cmp v0.5.6
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/kiegroup/kogito-operator/api v1.8.0
	github.com/kiegroup/kogito-operator/client v1.8.0
	github.com/stretchr/testify v1.7.0
	go.uber.org/zap v1.19.0
	k8s.io/api v0.21.4
	k8s.io/apimachinery v0.21.4
	k8s.io/client-go v0.21.4
	knative.dev/eventing v0.25.1-0.20210909163359-316e14d7fbc2
	knative.dev/hack v0.0.0-20210806075220-815cd312d65c
	knative.dev/pkg v0.0.0-20210909165259-d4505c660535
)

replace (
	github.com/kiegroup/kogito-operator/api => github.com/kiegroup/kogito-operator/api v0.0.0-20210702132500-6452df3eb8be
	github.com/kiegroup/kogito-operator/client => github.com/kiegroup/kogito-operator/client v0.0.0-20210702132500-6452df3eb8be
	github.com/prometheus/client_golang => github.com/prometheus/client_golang v0.9.2
)
