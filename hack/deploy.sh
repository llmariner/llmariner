#!/usr/bin/env bash

set -euo pipefail

cluster_name="llm-operator-demo"

kind load docker-image inference-server:latest -n "${cluster_name}"
kind load docker-image job-manager-server:latest -n "${cluster_name}"

kubectl apply -f inference-server.yaml
kubectl apply -f job-manager.yaml
