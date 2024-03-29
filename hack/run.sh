#!/usr/bin/env bash

set -euo pipefail

cluster_name="llm-operator-demo"

kind create cluster --name "${cluster_name}"

kind load docker-image inference-server:latest -n "${cluster_name}"
kind load docker-image job-manager-server:latest -n "${cluster_name}"

kubectl apply -f inference-server.yaml
kubectl apply -f job-manager.yaml

kubectl port-forward service/inference-server 11434:11434 &
kubectl port-forward service/job-manager-server 8080:8080 &

# Send a test request.
curl http://localhost:11434/api/generate -d '{
  "model": "gemma:2b",
  "prompt":"Why is the sky blue?"
}'

# Test gRPC enpdoint.
grpcurl -plaintext localhost:8080 list llmoperator.job_manager.server.v1.JobManagerService
