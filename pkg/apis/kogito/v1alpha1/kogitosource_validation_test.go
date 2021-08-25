/*
Copyright 2020 The Knative Authors

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

package v1alpha1

import (
	"context"
	"testing"

	"github.com/kiegroup/kogito-operator/api/v1beta1"
	duckv1 "knative.dev/pkg/apis/duck/v1"

	"github.com/google/go-cmp/cmp"
	"knative.dev/pkg/webhook/resourcesemantics"

	"knative.dev/pkg/apis"
)

func TestKogitoSourceValidation(t *testing.T) {
	testCases := map[string]struct {
		cr   resourcesemantics.GenericCRD
		want *apis.FieldError
	}{
		"nil spec": {
			cr: &KogitoSource{
				Spec: KogitoSourceSpec{
					KogitoRuntimeSpec: v1beta1.KogitoRuntimeSpec{
						KogitoServiceSpec: v1beta1.KogitoServiceSpec{
							Image: "quay.io/kiegroup/serverless-workflow-example",
						},
					},
				},
			},
			want: func() *apis.FieldError {
				var errs *apis.FieldError

				feSink := apis.ErrGeneric("expected at least one, got none", "ref", "uri")
				feSink = feSink.ViaField("sink").ViaField("spec")
				errs = errs.Also(feSink)

				feServiceAccountName := apis.ErrMissingField("serviceAccountName")
				feServiceAccountName = feServiceAccountName.ViaField("spec")
				errs = errs.Also(feServiceAccountName)

				return errs
			}(),
		},

		"no kogito image": {
			cr: &KogitoSource{
				Spec: KogitoSourceSpec{
					SourceSpec: duckv1.SourceSpec{
						Sink: duckv1.Destination{
							Ref: &duckv1.KReference{
								Kind:       "Pod",
								Namespace:  "default",
								Name:       "mypod",
								APIVersion: "v1",
							},
						},
					},
					ServiceAccountName: "default",
					KogitoRuntimeSpec:  v1beta1.KogitoRuntimeSpec{},
				},
			},
			want: func() *apis.FieldError {
				var errs *apis.FieldError

				feImage := apis.ErrMissingField("image")
				feImage = feImage.ViaField("spec")
				errs = errs.Also(feImage)

				return errs
			}(),
		},
	}

	for n, test := range testCases {
		t.Run(n, func(t *testing.T) {
			got := test.cr.Validate(context.Background())
			if diff := cmp.Diff(test.want.Error(), got.Error()); diff != "" {
				t.Errorf("%s: validate (-want, +got) = %v", n, diff)
			}
		})
	}
}
