#!/usr/bin/env bash

# Copyright 2021 The Knative Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


script_dir_path=`dirname "${BASH_SOURCE[0]}"`

set -e
source ${script_dir_path}/../ci/operator-ensure-manifests.sh

# SCRIPT_URL URL to the script used by OLM to test the operator
SCRIPT_URL="https://raw.githubusercontent.com/redhat-openshift-ecosystem/operator-test-playbooks/master/upstream/test/test.sh"

cd "${tempfolder}"

bash <(curl -sL "${SCRIPT_URL}") all  community-operators/eventing-kogito/"${version}"
