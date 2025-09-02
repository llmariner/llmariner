#!/usr/bin/env bash

set -euo pipefail
trap 'kill $(jobs -p)' EXIT

echo "Creating a base model..."

llma models create base TinyLlama/TinyLlama-1.1B-Chat-v1.0 --source-repository hugging-face

base_model_id=TinyLlama-TinyLlama-1.1B-Chat-v1.0

for i in {1..300}; do
  if llma models list | grep "${base_model_id}" | grep succeeded; then
	break
  fi
  sleep 1
done

echo "Creating a fine-tuned model..."

# TODO(kenji): Use a proper model. Currently the test assumes that inference-sim is used as a runtime and
# actual adapter loading is not tested.
llma models create fine-tuned \
  --base-model-id "${base_model_id}" \
  --source-repository object-store \
  --model-file-location s3://llm-operator-models/v1/workspace/fake-adapter \
  --suffix test

model_id=ft:TinyLlama/TinyLlama-1.1B-Chat-v1.0:test

# Wait until the model is loaded. The status of the model becomes "succeeded" when it is loaded.
for i in {1..300}; do
  if llma models list | grep "${model_id}"| grep succeeded; then
	break
  fi
  sleep 1
done

llma models activate "${model_id}"

llma chat completions create --model "${model_id}" --role user --completion "Hello"
