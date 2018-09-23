set -o errexit
set -o nounset
set -o pipefail

# Unset CDPATH so that path interpolation can work correctly
# https://github.com/bfsrnetes/bfsrnetes/issues/52255
unset CDPATH

# The root of the build/dist directory
if [ -z "$BFS_ROOT" ]
then   
    BFS_ROOT="$(cd "$(dirname "${BASH_SOURCE}")/../.." && pwd -P)"
if