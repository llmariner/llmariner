#!/usr/bin/env bash

set -euo pipefail

echo "Waiting for the deployment to be ready..."

kubectl wait --timeout=300s --for=condition=ready pod -n llmariner -l app.kubernetes.io/instance=llmariner

echo "Deployment is ready!"

echo "Waiting for google-gemma-2b-it-q4_0 to be loaded..."

for i in {1..300}; do
  if llma models list | grep google-gemma-2b-it-q4_0; then
    break
  fi
  sleep 1
done

echo "Model google-gemma-2b-it-q4_0 is loaded!"

echo "Running chat completion..."

llma chat completions create --model google-gemma-2b-it-q4_0 --role user --completion  "What is the capital of France?"

echo "Waiting for TinyLlama-TinyLlama-1.1B-Chat-v1.0-awq to be loaded..."

for i in {1..300}; do
  if llma models list | grep TinyLlama-TinyLlama-1.1B-Chat-v1.0-awq; then
    break
  fi
  sleep 1
done

echo "Model TinyLlama-TinyLlama-1.1B-Chat-v1.0-awq is loaded!"

llma chat completions create --model TinyLlama-TinyLlama-1.1B-Chat-v1.0-awq --role user --completion  "What is the capital of France?"

echo "Chat completion is done!"

# TODO(kenji): Test more.
