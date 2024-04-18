#!/usr/bin/env bash
#
# Deploy https://github.com/run-ai/fake-gpu-operator
#
# This is useful to simulate nvidia GPU resources locally.

set -euo pipefail

kubectl get nodes -o name | xargs -I {} kubectl label {} nvidia.com/gpu.deploy.device-plugin=true nvidia.com/gpu.deploy.dcgm-exporter=true --overwrite

helm repo add fake-gpu-operator https://fake-gpu-operator.storage.googleapis.com
helm repo update
helm upgrade -i gpu-operator fake-gpu-operator/fake-gpu-operator --namespace gpu-operator --create-namespace
