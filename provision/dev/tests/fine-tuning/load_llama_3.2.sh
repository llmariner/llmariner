#!/usr/bin/env bash

set -euo pipefail

echo "Creating a model..."
llma models create base meta-llama/Llama-3.2-1B-Instruct -s object-store

# Wait until the model is loaded. The status of the model becomes "succeeded" when it is loaded.
for i in {1..300}; do
  if llma models list | grep meta-llama-Llama-3.2-1B-Instruct | grep succeeded; then
	break
  fi
  sleep 1
done
