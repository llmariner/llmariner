import os
import sys

from openai import OpenAI

if len(sys.argv) != 2:
    print("Usage: python submit_fine_tuning_job.py <training_file>")
    sys.exit(1)

training_file = sys.argv[1]

client = OpenAI(
    base_url="https://api.llm.staging.cloudnatix.com/v1",
    api_key=os.getenv("LLMARINER_API_KEY"),
)

resp = client.fine_tuning.jobs.create(
    model="meta-llama-Llama-3.2-1B-Instruct",
    suffix='fine-tuning',
    training_file=training_file,
    hyperparameters={
        "n_epochs": 20,
    }
)
print('Created job. ID=%s' % resp.id)
