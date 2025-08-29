#!/usr/bin/env bash

set -euo pipefail

echo "Creating a model..."
llma models create base TheBloke/TinyLlama-1.1B-Chat-v1.0-GGUF/tinyllama-1.1b-chat-v1.0.Q2_K.gguf --source-repository hugging-face

# Wait until the model is loaded. The status of the model becomes "succeeded" when it is loaded.
for i in {1..300}; do
  if llma models list | grep TheBloke-TinyLlama-1.1B-Chat-v1.0-GGUF | grep succeeded; then
	break
  fi
  sleep 1
done

echo "Model is loaded!"

llma models get TheBloke-TinyLlama-1.1B-Chat-v1.0-GGUF-tinyllama-1.1b-chat-v1.0.Q2_K.gguf

echo "Activating the model."

llma models activate TheBloke-TinyLlama-1.1B-Chat-v1.0-GGUF-tinyllama-1.1b-chat-v1.0.Q2_K.gguf

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

llma chat completions create --model TheBloke-TinyLlama-1.1B-Chat-v1.0-GGUF-tinyllama-1.1b-chat-v1.0.Q2_K.gguf  --role user --completion "What is the capital of France?"

echo "Chat completion is done!"
