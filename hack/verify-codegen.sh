#!/usr/bin/env bash
#
# Copyright 2026 The Nephio Authors.
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
#
# Compare the committed API artifacts against what controller-gen produces.
#
# Both sides are snapshots under a temporary directory, so this never writes to
# the work tree: it is safe over a dirty checkout, and a failure cannot damage
# anything. Keep it that way. An earlier version created the CRD output
# directory in the work tree to make the comparison line up, which passed
# locally and failed on every clean checkout, because git does not track the
# empty directory this repository is supposed to have.

set -o errexit
set -o nounset
set -o pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

# Absolute, because the generators run after cd into a snapshot where a
# relative path would point at nothing: bin/ is not copied.
CONTROLLER_GEN="${CONTROLLER_GEN:-${REPO_ROOT}/bin/controller-gen}"
if [[ "${CONTROLLER_GEN}" != /* ]]; then
    CONTROLLER_GEN="$(cd "$(dirname "${CONTROLLER_GEN}")" && pwd)/$(basename "${CONTROLLER_GEN}")"
fi

# What this check owns, and nothing else. config/rbac is deliberately excluded:
# the committed role still grants randeployments, a resource this operator no
# longer has, while the ClusterRole actually deployed lives in the catalog.
# Verifying it here would freeze a stale file as the canonical output.
VERIFIED_YAML_DIR="config/crd/bases"
VERIFIED_GO_DIR="api"
GENERATED_GO_GLOB="zz_generated*.go"

if [[ ! -x "${CONTROLLER_GEN}" ]]; then
    echo "controller-gen not found at ${CONTROLLER_GEN}. Run 'make controller-gen'." >&2
    exit 1
fi

WORK_DIR="$(mktemp -d)"
trap 'rm -rf "${WORK_DIR}"' EXIT

EXPECTED_DIR="${WORK_DIR}/expected"
GENERATED_DIR="${WORK_DIR}/generated"
mkdir -p "${EXPECTED_DIR}" "${GENERATED_DIR}"

tar -c --exclude=./.git --exclude=./bin -C "${REPO_ROOT}" . | tar -x -C "${EXPECTED_DIR}"
tar -c -C "${EXPECTED_DIR}" . | tar -x -C "${GENERATED_DIR}"

# An absent output directory and an empty one describe the same committed
# artifact set, and this repository is meant to have no CRDs at all. Normalize
# both snapshots so the comparison is about files rather than about whether git
# happened to carry a directory.
mkdir -p "${EXPECTED_DIR}/${VERIFIED_YAML_DIR}"

# controller-gen only writes files, it never removes one it has stopped
# emitting, so the outputs are cleared in the generated snapshot first.
# Without this an artifact that should have disappeared survives every check.
rm -rf "${GENERATED_DIR:?}/${VERIFIED_YAML_DIR}"
mkdir -p "${GENERATED_DIR}/${VERIFIED_YAML_DIR}"
find "${GENERATED_DIR}/${VERIFIED_GO_DIR}" -name "${GENERATED_GO_GLOB}" -delete

# Scoped to ./api/... so that what is produced is exactly what is compared.
# make update-api-codegen runs the same two commands against the work tree.
(
    cd "${GENERATED_DIR}"
    "${CONTROLLER_GEN}" object:headerFile="hack/boilerplate.go.txt" paths="./api/..."
    "${CONTROLLER_GEN}" crd paths="./api/..." \
        output:crd:artifacts:config="${VERIFIED_YAML_DIR}"
)

status=0

# diff exits 1 when the trees differ and 2 or more when the comparison itself
# failed. Collapsing the two would report a broken tool as stale code
# generation and send someone off to regenerate files that were never wrong.
#
# -ru rather than -ruN: -N treats an absent file as an empty one, so an extra
# or missing zero-length artifact would compare equal and the inventory would
# not be an inventory. Both roots are normalized above, so -N is not needed.
compare_tree() {
    local left="$1"
    local right="$2"
    local result=0

    diff -ru "${left}" "${right}" || result=$?

    case "${result}" in
        0) ;;
        1) status=1 ;;
        *)
            printf 'unable to compare %s and %s (diff exited %d)\n' \
                "${left}" "${right}" "${result}" >&2
            exit "${result}"
            ;;
    esac
}

compare_tree "${EXPECTED_DIR}/${VERIFIED_YAML_DIR}" "${GENERATED_DIR}/${VERIFIED_YAML_DIR}"
compare_tree "${EXPECTED_DIR}/${VERIFIED_GO_DIR}" "${GENERATED_DIR}/${VERIFIED_GO_DIR}"

if [[ "${status}" -ne 0 ]]; then
    cat >&2 <<'EOF'

Committed generated artifacts do not match controller-gen output.
Run 'make update-api-codegen' and commit the result.
EOF
    exit 1
fi

echo "Generated artifacts are up to date."
