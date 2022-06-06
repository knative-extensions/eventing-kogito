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

set -e

script_dir_path=`dirname "${BASH_SOURCE[0]}"`
source ${script_dir_path}/../env.sh

tempfolder=$(mktemp -d)
echo "Temporary folder is ${tempfolder}"

version=$(getOperatorVersion)

git clone https://github.com/k8s-operatorhub/community-operators.git "${tempfolder}"
mv "${tempfolder}/operators" "${tempfolder}/community-operators"
mkdir  "${tempfolder}/community-operators/eventing-kogito/${version}/"
## copy the latest manifests
cp -r bundle/manifests/ "${tempfolder}/community-operators/eventing-kogito/${version}/"
cp -r bundle/metadata/ "${tempfolder}/community-operators/eventing-kogito/${version}/"
cp -r bundle/tests/ "${tempfolder}/community-operators/eventing-kogito/${version}/"
cp bundle.Dockerfile "${tempfolder}/community-operators/eventing-kogito/${version}/Dockerfile"

#Edit dockerfile with correct relative path
sed -i "s|bundle/manifests|manifests|g" "${tempfolder}/community-operators/eventing-kogito/${version}/Dockerfile"
sed -i "s|bundle/metadata|metadata|g" "${tempfolder}/community-operators/eventing-kogito/${version}/Dockerfile"
sed -i "s|bundle/tests|tests|g" "${tempfolder}/community-operators/eventing-kogito/${version}/Dockerfile"
