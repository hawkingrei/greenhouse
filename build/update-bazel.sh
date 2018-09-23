#!/usr/bin/env bash
# Copyright 2016 The Kubernetes Authors.
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

set -o errexit
set -o pipefail

export BFS_ROOT=$(dirname "${BASH_SOURCE}")/..
source "${BFS_ROOT}/build/lib/init.sh"

#bfs::util::ensure-gnu-sed

# Remove generated files prior to running kazel.
# TODO(spxtr): Remove this line once Bazel is the only way to build.
# rm -f "${BFS_ROOT}/pkg/generated/openapi/zz_generated.openapi.go"

# Ensure that we find the binaries we build before anything else.
export GOBIN="${BFS_OUTPUT_BINPATH}"
PATH="${GOBIN}:${PATH}"

# Install tools we need, but only from vendor/...
go install ./vendor/github.com/hawkingrei/kazel

# gazelle gets confused by our staging/ directory, prepending an extra
# "k8s.io/bfsrnetes/staging/src" to the import path.
# gazelle won't follow the symlinks in vendor/, so we can't just exclude
# staging/. Instead we just fix the bad paths with sed.

kazel
