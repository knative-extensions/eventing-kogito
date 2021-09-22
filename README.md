# Knative Eventing Kogito

[![GoDoc](https://godoc.org/knative-sandbox/eventing-kogito?status.svg)](https://godoc.org/knative-sandbox/eventing-kogito)
[![Go Report Card](https://goreportcard.com/badge/knative-sandbox/eventing-kogito)](https://goreportcard.com/report/knative-sandbox/eventing-kogito)
[![Releases](https://img.shields.io/github/release-pre/knative-sandbox/eventing-kogito.svg)](https://github.com/knative-sandbox/eventing-kogito/releases)
[![LICENSE](https://img.shields.io/github/license/knative-sandbox/eventing-kogito.svg)](https://github.com/knative-sandbox/eventing-kogito/blob/master/LICENSE)
[![codecov](https://codecov.io/gh/knative-sandbox/eventing-kogito/branch/main/graph/badge.svg)](https://codecov.io/gh/knative-sandbox/eventing-kogito)
[![zulip chat](https://img.shields.io/badge/zulip-join_chat-brightgreen.svg?logo=zulip)](https://kie.zulipchat.com/#narrow/stream/262892-serverless-workflow)
[![Slack](https://img.shields.io/badge/%23eventing-white.svg?logo=slack&color=522a5e)](https://knative.slack.com/archives/CQBKVH4QY)


[Kogito](https://kogito.kie.org/) is a platform to build cloud-native business automation services. It has a built-in
engine able to
run [BPMN](https://docs.jboss.org/kogito/release/latest/html_single/#chap-kogito-developing-process-services),
[Rules](https://docs.jboss.org/kogito/release/latest/html_single/#chap-kogito-using-drl-rules),
[DMN](https://docs.jboss.org/kogito/release/latest/html_single/#collection-kogito-developing-decision-services), and
[Serverless Workflows](https://docs.jboss.org/kogito/release/latest/html_single/#chap-kogito-orchestrating-serverless)
services.

The Knative Eventing Kogito is a way to have
your [Kogito Services in the Knative Eventing platform](https://docs.jboss.org/kogito/release/latest/html_single/#con-knative-eventing_kogito-developing-process-services).

This project is under active development. To know more about it, please see our [roadmap](./ROADMAP.md).

To learn more about Knative, please visit the
[Knative docs](https://github.com/knative/docs) repository.

If you are interested in contributing, see [CONTRIBUTING.md](./CONTRIBUTING.md)
and [DEVELOPMENT.md](./DEVELOPMENT.md).

## Kogito Source

You can deploy a Kogito service to act as a source of events, meaning that your business automation service can produce
events targeting any addressable service within the platform.

### Published event types and attributes

Events published by the Kogito Source are highly related to the type of service you create with Kogito. Usually, the
event that the service produces is tied to the business application domain.

#### Serverless Workflows and BPMN Process

If you have a [Serverless Workflow](https://github.com/serverlessworkflow/specification) definition that produces an
"Order" event, your data will reflect that structure. See the excerpt below:

```yaml
id: shippinghandling
name: Shipping Handling
start: ShippingHandling
version: "1.0"
events:
  - kind: produced
    name: InternationalShippingOrder
    type: internationalShipping
    source: internationalShipping
  - kind: produced
    name: DomesticShippingOrder
    type: domesticShipping
    source: domesticShipping
states:
  - name: ShippingHandling
    type: switch
    dataConditions:
      - condition: "{{ $.[?(@.country == 'US')] }}"
        transition: DomesticShipping
      - condition: "{{ $.[?(@.country != 'US')] }}"
        transition: InternationalShipping
  - name: DomesticShipping
    type: inject
    data:
      shipping: "domestic"
    end:
      produceEvents:
        - eventRef: DomesticShippingOrder
  - name: InternationalShipping
    type: inject
    data:
      shipping: "international"
    end:
      produceEvents:
        - eventRef: "InternationalShippingOrder"
```

In this example, the Kogito Serverless Workflow service will produce either a `DomesticShippingOrder` or an
`InternationalShippingOrder`, depending on the data processed by the workflow.

The developer of the workflow defines the data structure of the event. In this example, you can check the domain model
[in this repository](https://github.com/kiegroup/kogito-examples/tree/stable/serverless-workflow-order-processing). The
domain is based on the order model sent by the [main workflow](https://github.com/kiegroup/kogito-examples/blob/stable/serverless-workflow-order-processing/src/main/resources/order-workflow.sw.yaml).

An example of a `CloudEvent` produced by this service looks like this one:

```log
cloudevents.Event
Validation: valid
Context Attributes,
  specversion: 1.0
  type: internationalShipping
  source: /process/shippinghandling
  id: d557fad8-81b5-482a-b981-5ecb267a92f9
  time: 2021-04-13T20:40:38.431677Z
Extensions,
  kogitoparentprociid: f12e91c0-8980-40b6-a49b-3c35ce435718
  kogitoprocid: shippinghandling
  kogitoprocinstanceid: ba57743d-521f-41f2-864a-a7b3f68d35af
  kogitorootprocid: orderworkflow
  kogitorootprociid: f12e91c0-8980-40b6-a49b-3c35ce435718
  kogitousertaskist: 1
Data,
  {"id":"f0643c68-609c-48aa-a820-5df423fa4fe0","country":"Brazil","total":10000,"description":"iPhone 12","fraudEvaluation":true,"shipping":"international"}
```

##### CloudEvents Structure

This is the structure of a produced `CloudEvent`:

| Attribute         | Value                     | Notes | 
| ----------------- | ------------------------- | ----- |
| `type`            | **Serverless Workflow:** it takes the type defined in the workflow [event definition](https://github.com/serverlessworkflow/specification/blob/main/specification.md#event-definition) <br/>**BPMN:** uses the message trigger name defined in the process | | 
| `source`          | BPMN and Serverless Workflow have  a constant `/process/` followed by the BPMN process id or workflow id, respectively. | |
| `id`              | An unique generated ID | |
| `data`            | Contains the domain model of the event produced by the service. It depends on the business model defined by the developer. | |
| `extensions`      | CloudEvent extensions that carry Kogito engine data related to the workflow instance, such as the instance id of the workflow. | All extensions produced by the service start with `kogito`. See the table below for a list of all extensions that the Kogito Eventing supports. |
| `datacontenttype` | `application/cloudevents+json` | See the [CloudEvents spec reference](https://github.com/cloudevents/spec/blob/v1.0.1/json-format.md#3-envelope) for this media type |

Following the possible `CloudEvents` extensions that can be added to the produced events:

| Extension              | Name                        | Notes                                                   |
| ---------------------- | --------------------------- | ------------------------------------------------------- |
| `kogitoprocinstanceid` | Process Instance ID         | Generated ID for a Serverless Workflow or BPMN instance |
| `kogitoprocrefid`      | Process Reference ID        |  |
| `kogitoprocist`        | Process Instance State      |  |
| `kogitoprocid`         | Process ID                  | Identification of the process. For Serverless Workflows, it's the "workflow ID" attribute in the definition |
| `kogitoparentprociid`  | Parent Process Instance ID  | Identification of the parent process or workflow |
| `kogitorootprociid`    | Root Process Instance ID    | Identification of the root instance process or workflow |
| `kogitorootprocid`     | Root Process ID             | Identification of the root process or workflow. [Subflows](https://github.com/serverlessworkflow/specification/blob/main/specification.md#subflow-action) will have this attribute set to the root workflow ID. |
| `kogitoprocstartfrom`  | Process Start from Node     | ID of the node where the process or workflow started from. |
| `kogitousertaskiid`    | User Task Instance ID       | Identification for BPMN User Task Instances. |
| `kogitousertaskist`    | User Task Instance State    | ID of the User Task instance. |
| `kogitoaddon`          | Kogito Add-ons              | Collection of [Kogito Add-ons](https://github.com/kiegroup/kogito-runtimes/tree/main/addons) used by this process or workflow instance. |

Please visit our documentation to know more about creating your Kogito BPMN or Serverless Workflow with Knative
Eventing.

### Samples

In the [`examples`](./examples) directory, you will find some examples to get started with Kogito Eventing. To name a
few:

- [**Order processing workflow**](./examples): a small service capable of consuming an event and produce different
  outputs based on the contents of the order domain
- [**Camel Telegram**](./examples/camel-telegram): Redirect your events to a Telegram bot chat

### Installation

How to install the Kogito Source in your cluster.

#### Prerequisites

Before installing the Knative Eventing Kogito Source, you must meet the following prerequisites:

1. You have administrative privileges in the target cluster
2. You have [installed Knative](https://knative.dev/docs/install/) Eventing and Serving (or
   have [OpenShift Serverless Platform](https://www.openshift.com/learn/topics/serverless) available)
3. You have [installed the Kogito Operator](https://github.com/kiegroup/kogito-operator)

#### Installation steps

You can install the source using `kubectl` CLI:

```shell
VERSION=0.26.0
kubectl apply -f https://github.com/knative-sandbox/eventing-kogito/releases/download/v${VERSION}/kogito.yaml 
```

Replace the `VERSION` variable with the desired target version.

> Note that Kogito Eventing is only available from version **0.26.0**.

By running the above command you will create the namespace `knative-kogito` in your cluster and all the resources
necessary to run the Kogito Source:

```shell
kubectl get pods -n knative-kogito

NAME                                        READY   STATUS    RESTARTS   AGE
kogito-source-controller-7689d9dc6d-4l5hv   2/2     Running   1          86m
kogito-source-webhook-54b7c87ff9-5r46h      1/1     Running   0          127m
```

Now you can start deploying the [examples](./examples)!

### Using the event source

#### Deploying the Kogito Event Source

1. Deploy a [Knative Sink](https://knative.dev/docs/developer/eventing/sinks/) to consume the events produced by your
   Kogito Service
2. Create your Kogito project locally. You can use our guide
   for [BPMN](https://docs.jboss.org/kogito/release/latest/html_single/#con-knative-eventing_kogito-developing-process-services)
   or [Serverless Workflow](https://docs.jboss.org/kogito/release/latest/html_single/#chap-kogito-orchestrating-serverless)
   .
3. Build the image with your project. See
   an [example here](https://github.com/kiegroup/kogito-examples/blob/stable/serverless-workflow-order-processing/Dockerfile)
   .
4. Push your image to a registry where your cluster can access it.
5. Create the Kogito Source CR. See [this example](./examples/kogito-source-reference.yaml) to be used as a reference to
   create your own.
6. Use kubectl to deploy your source: `kubectl apply -f <path to the source>.yaml`

#### Verifying deployed sources

Having deployed your sources you can easily list them and check their status with the following command:

```shell
kubectl get kogitosources
```

You should see a listing like this one:

```shell
NAME                      READY   REASON   SINK                                                        AGE
kogito-order-processing   True             http://kogito-channel-kn-channel.kogito.svc.cluster.local   31s
```

The `READY` field gives a cue about the general status of the source.

Behind the curtains, Knative Eventing Kogito creates a
backed [`KogitoRuntime`](https://github.com/kiegroup/kogito-operator/blob/main/apis/app/v1beta1/kogitoservices.go)
object that is a representation of your service in the cluster. You can also query it to check its status:

```shell
kubectl get kogitoruntimes

NAME                         REPLICAS   IMAGE                                                    ENDPOINT
ks-kogito-order-processing   1          quay.io/ricardozanini/order-processing-workflow:latest   
```

The source is ready to produce events to the sink defined in the service.

#### Objects deployed by the source

The image below illustrates the objects created and managed by Kogito Knative Source:

![Knative Kogito Source](./docs/knative-kogito-source-role.png)

Once you deploy the `KogitoSource` object, the source controller will discover the endpoint of the addressable object
defined in the `sink` attribute.

Having resolved the endpoint address, the controller then delegates the creation of the Kogito service to the Kogito
Operator. The operator controller handles the configuration of the sink endpoint in the service container.

## Additional Resources

- [Kogito Official Documentation](https://docs.jboss.org/kogito/release/latest/html_single/)
- [Blog: Orchestrating Events with Knative and Kogito](https://knative.dev/blog/2020/12/17/orchestrating-events-with-knative-and-kogito/)

## Roadmap

To learn more about future features, check our [Roadmap](ROADMAP.md).
