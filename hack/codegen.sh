#!/bin/bash -e

# Copyright 2016 The Rook Authors. All rights reserved.
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

#scriptdir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

#cd ${scriptdir}/../vendor/k8s.io/code-generator && ./generate-groups.sh \
#  all \
#  practice-operator/pkg/client \
#  practice-operator/pkg/apis \
#  "student:v1" 
#  #--go-header-file ${scriptdir}/../hack/boilerplate/boilerplate.go.txt




set -o errexit
set -o nounset
set -o pipefail

#SCRIPT_ROOT=$(dirname ${BASH_SOURCE[0]})/..
SCRIPT_ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
CODEGEN_PKG=${CODEGEN_PKG:-$(cd ${SCRIPT_ROOT}; ls -d -1 ../vendor/k8s.io/code-generator 2>/dev/null || echo ../code-generator)}
${CODEGEN_PKG}/generate-groups.sh all \
        practice-operator/pkg/client \
        practice-operator/pkg/apis \
        "student:v1"
       # --go-header-file ${SCRIPT_ROOT}/boilerplate/boilerplate.go.txt

       # --output-base "$(dirname ${BASH_SOURCE})/../../../../" \
