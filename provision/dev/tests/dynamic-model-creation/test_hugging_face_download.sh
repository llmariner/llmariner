#!/usr/bin/env bash

set -euo pipefail

echo "Creating a model..."
llma models create QuantFactory/SmolLM-135M-GGUF/SmolLM-135M.Q2_K.gguf -s hugging-face

# Wait until the model is loaded. The status of the model becomes "succeeded" when it is loaded.
for i in {1..300}; do
  if llma models list | grep QuantFactory-SmolLM-135M-GGUF-SmolLM-135M.Q2_K.gguf | grep succeeded; then
	break
  fi
  sleep 1
done

echo "Model is loaded!"
echo "Running chat completion..."
llma chat completions create --model QuantFactory-SmolLM-135M-GGUF-SmolLM-135M.Q2_K.gguf  --role user --completion "What is the capital of France?"
