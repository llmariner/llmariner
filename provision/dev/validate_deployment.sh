#!/usr/bin/env bash

set -euo pipefail

echo "Waiting for the deployment to be ready..."

kubectl wait --timeout=300s --for=condition=ready pod -n llmariner -l app.kubernetes.io/instance=llmariner

echo "Deployment is ready!"

echo "Waiting for the model to be loaded..."

for i in {1..300}; do
  if llma models list | grep google-gemma-2b-it-q4_0; then
    break
  fi
  sleep 1
done

echo "Model is loaded!"

echo "Waiting for the inference runtime pod is created..."

for i in {1..300}; do
  if kubectl get pod -n llmariner ollama-google-gemma-2b-it-q4-0-0; then
    break
  fi
  sleep 1
done

kubectl wait --timeout=300s --for=condition=ready pod -n llmariner ollama-google-gemma-2b-it-q4-0-0

echo "Inference runtime pod is ready!"

echo "Running chat completion..."

llma chat completions create --model google-gemma-2b-it-q4_0 --role user --completion  "What is the capital of France?"

echo "Chat completion is done!"

# TODO(kenji): Test more.
