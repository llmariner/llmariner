#!/usr/bin/env bash

set -euo pipefail

cluster_name="kind-gpu"

./nvkind cluster create --name "${cluster_name}" --config-template ./kind-cluster.yaml --config-values=- \
<<EOF
numGPUs: 1
EOF
