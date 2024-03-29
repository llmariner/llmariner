#!/usr/bin/env bash

set -euo pipefail

cluster_name="llm-operator-demo"

kind create cluster --name "${cluster_name}"

kind load docker-image inference-server:latest -n "${cluster_name}"

kubectl apply -f inference-server.yaml

kubectl port-forward service/inference-server 11434:11434 &

# Send a test request.
curl http://localhost:11434/api/generate -d '{
  "model": "gemma:2b",
  "prompt":"Why is the sky blue?"
}'
