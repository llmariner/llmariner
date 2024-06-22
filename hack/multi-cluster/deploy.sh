#!/usr/bin/env bash

set -euo pipefail

basedir=$(dirname "$0")

# Set the kubectl context to the control plane cluster.
kubectl config use-context kind-llm-operator-control-plane
"${basedir}"/deploy_control_plane.sh

llmo auth login

# Set the kubectl context to the worker plane cluster.
kubectl config use-context kind-llm-operator-worker-plane
"${basedir}"/deploy_worker_plane.sh
