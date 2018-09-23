#!/bin/bash
set -o errexit
set -o nounset
set -o pipefail

export BFS_ROOT=$(dirname "${BASH_SOURCE}")/..
source "${BFS_ROOT}/build/lib/init.sh"

if [ $(uname -s) = "Linux" ]; then 
        bfs::util::ensure-bazel 
fi

if [ $(uname -s) = "Darwin" ]; 
then 
        bfs::util::ensure-homebrew 
        bfs::util::ensure-homebrew-bazel 
fi

