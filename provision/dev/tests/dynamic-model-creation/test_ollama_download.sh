#!/usr/bin/env bash

set -euo pipefail

echo "Creating a model..."
llma models create deepseek-r1:1.5b -s ollama

# Wait until the model is loaded. The status of the model becomes "succeeded" when it is loaded.
for i in {1..300}; do
  if llma models list | grep deepseek-r1:1.5b | grep succeeded; then
	break
  fi
  sleep 1
done

echo "Model is loaded!"
echo "Running chat completion..."
llma chat completions create base --model deepseek-r1:1.5b --role user --completion "What is the capital of France?"
