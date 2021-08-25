/*
Copyright 2021 The Knative Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package reconciler

import (
	"bytes"
	"context"
	"io/ioutil"
	"testing"

	kogitoclient "github.com/kiegroup/kogito-operator/client/clientset/versioned"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"

	kogitofake "github.com/kiegroup/kogito-operator/client/clientset/versioned/fake"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
	k8sfake "k8s.io/client-go/kubernetes/fake"
	"knative.dev/eventing-kogito/pkg/apis/kogito/v1alpha1"
	"knative.dev/eventing-kogito/pkg/reconciler/kogito/resources"
	v1 "knative.dev/eventing/pkg/apis/sources/v1"
	"knative.dev/pkg/reconciler"
)

type exampleKogitoSource string

const (
	orderProcessingWorkflow exampleKogitoSource = "../../examples/order-processing-workflow.yaml"
)

func getExampleKogitoSource(example exampleKogitoSource) *v1alpha1.KogitoSource {
	kogitoSourceBytes, err := ioutil.ReadFile(string(example))
	if err != nil {
		panic(err)
	}
	src := &v1alpha1.KogitoSource{}
	decoder := yaml.NewYAMLOrJSONDecoder(bytes.NewReader(kogitoSourceBytes), 1000)
	if err := decoder.Decode(src); err != nil {
		panic(err)
	}
	return src
}

type reconcileHandler struct {
	KubeClientSet   kubernetes.Interface
	KogitoClientSet kogitoclient.Interface
}

func (r *reconcileHandler) withKubeObjects(objects ...runtime.Object) *reconcileHandler {
	r.KubeClientSet = k8sfake.NewSimpleClientset(objects...)
	return r
}

func (r *reconcileHandler) withKogitoObjects(objects ...runtime.Object) *reconcileHandler {
	r.KogitoClientSet = kogitofake.NewSimpleClientset(objects...)
	return r
}

func (r *reconcileHandler) reconcile(src *v1alpha1.KogitoSource) (*appsv1.Deployment, *v1.SinkBinding, reconciler.Event) {
	// TODO: add to the context mocks for SinkBinding Resolver
	ctx := context.TODO()
	if r.KogitoClientSet == nil {
		r.KogitoClientSet = kogitofake.NewSimpleClientset()
	}
	if r.KubeClientSet == nil {
		r.KubeClientSet = k8sfake.NewSimpleClientset()
	}
	runtimeReconciler := KogitoRuntimeReconciler{KubeClientSet: r.KubeClientSet, KogitoClientSet: r.KogitoClientSet}
	return runtimeReconciler.ReconcileKogitoRuntime(ctx, src,
		&v1.SinkBinding{
			ObjectMeta: metav1.ObjectMeta{
				Name:      src.GetName(),
				Namespace: src.GetNamespace(),
			},
			Spec: v1.SinkBindingSpec{
				SourceSpec: src.Spec.SourceSpec,
			},
		},
		resources.MakeReceiveAdapter(&resources.ReceiveAdapterArgs{
			EventSource: src.Namespace + "/" + src.Name,
			Source:      src,
			Labels:      resources.Labels(src.Name),
		}))
}

func newReconcileHandler() *reconcileHandler {
	return &reconcileHandler{}
}

func TestKogitoRuntimeReconciler_VerifyCreation(t *testing.T) {
	src := getExampleKogitoSource(orderProcessingWorkflow)
	deployment, binder, event := newReconcileHandler().reconcile(src)
	// we don't have Kogito Operator to handle the Deployment for us
	assert.Nil(t, deployment)
	assert.NotNil(t, binder)
	assert.True(t, reconciler.EventIs(event, newKogitoRuntimeCreated(src.Namespace, src.Name)))
}

func TestKogitoRuntimeReconciler_CheckChanges(t *testing.T) {
	src := getExampleKogitoSource(orderProcessingWorkflow)
	kogitoRuntime := resources.MakeReceiveAdapter(&resources.ReceiveAdapterArgs{
		EventSource: src.Namespace + "/" + src.Name,
		Source:      src,
		Labels:      resources.Labels(src.Name),
	})
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      kogitoRuntime.Name,
			Namespace: kogitoRuntime.Namespace,
		},
	}
	newImageName := "knative.dev/new/image:latest"
	newSrc := src.DeepCopy()
	newSrc.Spec.Image = newImageName
	// Our cluster already have a KogitoRuntime with the old name
	handler := newReconcileHandler().withKogitoObjects(kogitoRuntime).withKubeObjects(deployment)
	// Reconcile with the new image name
	_, binder, event := handler.reconcile(newSrc)
	assert.NotNil(t, binder)
	// the KogitoRuntime should be updated
	assert.True(t, reconciler.EventIs(event, newKogitoRuntimeUpdated(src.Namespace, src.Name)))
	updatedKogitoRuntime, err := handler.
		KogitoClientSet.AppV1beta1().KogitoRuntimes(src.Namespace).
		Get(context.TODO(), kogitoRuntime.Name, metav1.GetOptions{})
	assert.NoError(t, err)
	assert.Equal(t, newImageName, updatedKogitoRuntime.Spec.Image)
}
