#!/bin/bash

# Copyright 2014 The Kubernetes Authors.
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
set -o nounset
set -o pipefail

# Unset CDPATH so that path interpolation can work correctly
# https://github.com/bfsrnetes/bfsrnetes/issues/52255
unset CDPATH

# The root of the build/dist directory
BFS_ROOT="$(cd "$(dirname "${BASH_SOURCE}")/../.." && pwd -P)"

BFS_OUTPUT_SUBPATH="${BFS_OUTPUT_SUBPATH:-_output/local}"
BFS_OUTPUT="${BFS_ROOT}/${BFS_OUTPUT_SUBPATH}"
BFS_OUTPUT_BINPATH="${BFS_OUTPUT}/bin"

# This controls rsync compression. Set to a value > 0 to enable rsync
# compression for build container
BFS_RSYNC_COMPRESS="${BFS_RSYNC_COMPRESS:-0}"

# Set no_proxy for localhost if behind a proxy, otherwise,
# the connections to localhost in scripts will time out
export no_proxy=127.0.0.1,localhost
export bazel_version="0.15.2"
# This is a symlink to binaries for "this platform", e.g. build tools.
THIS_PLATFORM_BIN="${BFS_ROOT}/_output/bin"

source "${BFS_ROOT}/build/lib/util.sh"
source "${BFS_ROOT}/build/lib/logging.sh"

#bfs::log::install_errexit

#source "${BFS_ROOT}/build/lib/version.sh"
#source "${BFS_ROOT}/build/lib/golang.sh"
#source "${BFS_ROOT}/build/lib/etcd.sh"


# This emulates "readlink -f" which is not available on MacOS X.
# Test:
# T=/tmp/$$.$RANDOM
# mkdir $T
# touch $T/file
# mkdir $T/dir
# ln -s $T/file $T/linkfile
# ln -s $T/dir $T/linkdir
# function testone() {
#   X=$(readlink -f $1 2>&1)
#   Y=$(bfs::readlinkdashf $1 2>&1)
#   if [ "$X" != "$Y" ]; then
#     echo readlinkdashf $1: expected "$X", got "$Y"
#   fi
# }
# testone /
# testone /tmp
# testone $T
# testone $T/file
# testone $T/dir
# testone $T/linkfile
# testone $T/linkdir
# testone $T/nonexistant
# testone $T/linkdir/file
# testone $T/linkdir/dir
# testone $T/linkdir/linkfile
# testone $T/linkdir/linkdir
function bfs::readlinkdashf {
  # run in a subshell for simpler 'cd'
  (
    if [[ -d "$1" ]]; then # This also catch symlinks to dirs.
      cd "$1"
      pwd -P
    else
      cd "$(dirname "$1")"
      local f
      f=$(basename "$1")
      if [[ -L "$f" ]]; then
        readlink "$f"
      else
        echo "$(pwd -P)/${f}"
      fi
    fi
  )
}

# This emulates "realpath" which is not available on MacOS X
# Test:
# T=/tmp/$$.$RANDOM
# mkdir $T
# touch $T/file
# mkdir $T/dir
# ln -s $T/file $T/linkfile
# ln -s $T/dir $T/linkdir
# function testone() {
#   X=$(realpath $1 2>&1)
#   Y=$(bfs::realpath $1 2>&1)
#   if [ "$X" != "$Y" ]; then
#     echo realpath $1: expected "$X", got "$Y"
#   fi
# }
# testone /
# testone /tmp
# testone $T
# testone $T/file
# testone $T/dir
# testone $T/linkfile
# testone $T/linkdir
# testone $T/nonexistant
# testone $T/linkdir/file
# testone $T/linkdir/dir
# testone $T/linkdir/linkfile
# testone $T/linkdir/linkdir
bfs::realpath() {
  if [[ ! -e "$1" ]]; then
    echo "$1: No such file or directory" >&2
    return 1
  fi
  bfs::readlinkdashf "$1"
}

