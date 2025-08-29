#!/usr/bin/env bash

set -euo pipefail

echo "Creating a model..."
llma models create base deepseek-r1:1.5b --source-repository ollama

# Wait until the model is loaded. The status of the model becomes "succeeded" when it is loaded.
for i in {1..300}; do
  if llma models list | grep deepseek-r1:1.5b | grep succeeded; then
	break
  fi
  sleep 1
done

echo "Model is loaded!"

llma models get deepseek-r1:1.5b

echo "Waiting for the inference runtime pod is created..."

for i in {1..300}; do
  if kubectl get pod -n llmariner -l app.kubernetes.io/name=runtime 2>&1 | grep -v "No resources found"; then
    break
  fi
  sleep 1
done

kubectl wait --timeout=300s --for=condition=ready pod -n llmariner -l app.kubernetes.io/name=runtime

echo "Inference runtime pod is ready!"

echo "Running chat completion..."
llma chat completions create --model deepseek-r1:1.5b --role user --completion "What is the capital of France?"
