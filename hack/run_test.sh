#!/usr/bin/env bash

set -euo pipefail

# Send test requests

curl http://localhost:80/v1/models

curl http://localhost:80/v1/files

curl http://localhost:80/v1/chat/completions -d '{
  "model": "gemma:2b",
  "messages": [{"role": "user", "content": "Why is the sky blue?"}]
}'

# Test the fine-tuning service.
curl http://localhost:80/v1/fine_tuning/jobs

# Test OpenAI API by following https://platform.openai.com/docs/quickstart?context=python
python3 -m venv openai-env
source openai-env/bin/activate
pip3 install --upgrade openai

python3 run_completion.py
python3 run_fine_tuning.py
