# Kogito Source Examples

In this directory you will find some examples of Kogito Source usage. It's assumed that you have successfully installed
Knative and Knative Kogito Source in your cluster.

## Order Processing Workflow

This is a pre-built application based on the
original [Order Processing example](https://github.com/kiegroup/kogito-examples/tree/stable/serverless-workflow-order-processing)
. Please check that example to understand how to interact with the application.

To deploy this example in your cluster, first deploy the event display sink:

```shell
$ kubectl apply -f https://github.com/knative-sandbox/eventing-kogito/blob/main/examples/sinks/event-display.yaml
```

Then you can deploy the Kogito application with:

```shell
$ kubectl apply -f https://github.com/knative-sandbox/eventing-kogito/blob/main/examples/order-processing-workflow.yaml
```

Next, deploy the Kogito Source that will bind the Sink and the Service together:

```shell
$ kubectl apply -f https://github.com/knative-sandbox/eventing-kogito/blob/main/examples/kogito-source-order-processing-workflow.yaml
```

Check if the source is ready:

```shell
$ kubectl get kogitosource

NAME                             READY   REASON   SINK                                            AGE
kogito-order-processing-source   True             http://event-display.kogito.svc.cluster.local   110m
```

You can send "order" events to the application endpoint. You should see them in the Event Display service:

```
$ kubectl logs -l serving.knative.dev/service=event-display -c user-container

# suppressed (..)
Extensions,
  kogitoparentprociid: d896d34b-9650-41bb-b52e-dcabb58caa93
  kogitoprocid: fraudhandling
  kogitoprocinstanceid: 8876ff92-306e-4d8a-b31f-c833186eaf19
  kogitorootprocid: orderworkflow
  kogitorootprociid: d896d34b-9650-41bb-b52e-dcabb58caa93
  kogitousertaskist: 1
Data,
  {"id":"f0643c68-609c-48aa-a820-5df423fa4fe0","country":"Italy","total":10000,"description":"iPhone 12","shipping":"international","fraudEvaluation":true}
```

## Telegram CloudEvents Notification Bot

Please visit the [camel-telegram](./camel-telegram) directory for more details.
