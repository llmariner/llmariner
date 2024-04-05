#!/usr/bin/env bash

set -euo pipefail

kubectl port-forward -n inference-server service/inference-server 11434:11434 &
kubectl port-forward -n job-manager service/job-manager-server-http 8080:8080 &
kubectl port-forward -n job-manager service/job-manager-server-grpc 8081:8081 &
kubectl port-forward -n inference-manager service/inference-manager-engine-ollama 8082:8080 &

# Send a test request.
curl http://localhost:8082/api/generate -d '{
  "model": "gemma:2b",
  "prompt":"Why is the sky blue?"
}'

# Test the fine-tuning service.
curl http://localhost:8080/v1/fine_tuning/jobs
curl -X POST http://localhost:8080/v1/fine_tuning/jobs
grpcurl -plaintext localhost:8081 list llmoperator.fine_tuning.server.v1.FineTuningService

# Test OpenAI API by following https://platform.openai.com/docs/quickstart?context=python
python3 -m venv openai-env
source openai-env/bin/activate
pip3 install --upgrade openai

python3 run_openai.py
