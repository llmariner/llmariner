#!/usr/bin/env bash

set -euo pipefail

LLM_ENDPOINT_URL="${LLM_ENDPOINT_URL:-http://localhost/v1}"

echo "Checking if API endpoints are reachable."
paths=(
  basemodels
  files
  fine_tuning/jobs
  models
)
for path in "${paths[@]}"; do
  curl --fail --silent -H "Authorization: Bearer ${LLM_OPERATOR_TOKEN}" "${LLM_ENDPOINT_URL}/${path}" > /dev/null
done

echo "Checking if base models are loaded."
curl --silent -H "Authorization: Bearer ${LLM_OPERATOR_TOKEN}" "${LLM_ENDPOINT_URL}/basemodels" | jq -e '.data | map(select(.id == "google-gemma-2b-it-q4")) | length == 1' > /dev/null

echo "Checking if chat completions work."
curl --request POST --fail --silent "${LLM_ENDPOINT_URL}/chat/completions" -d '{
  "model": "google-gemma-2b-it-q4",
  "messages": [{"role": "user", "content": "Why is the sky blue?"}]
}' > /dev/null


echo "Passed."
