#!/usr/bin/env bash
#
# Deploy https://github.com/run-ai/fake-gpu-operator
#
# This is useful to simulate nvidia GPU resources locally.

set -euo pipefail

kubectl get nodes -o name | xargs -I '{}' kubectl label '{}' nvidia.com/gpu.deploy.device-plugin=true nvidia.com/gpu.deploy.dcgm-exporter=true --overwrite

helm repo add fake-gpu-operator https://fake-gpu-operator.storage.googleapis.com
helm repo update
# Pinned to 0.0.51 since 0.0.53 didn't work (and there is no 0.0.52).
helm upgrade -i gpu-operator fake-gpu-operator/fake-gpu-operator --namespace gpu-operator --create-namespace --version 0.0.51
