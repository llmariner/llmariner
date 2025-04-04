import json
import os
import sys
import time

from openai import OpenAI
from datasets import load_dataset

client = OpenAI(
    # TODO(kenji): Be able to change this to http://localhost:80 for the multi-cluster
    # configuration.
    base_url="http://localhost:8080/v1",
    api_key=os.getenv("LLMARINER_API_KEY"),
)

dataset = load_dataset("csv", data_files={
    "train": "https://raw.githubusercontent.com/AlexandrosChrtn/llama-fine-tune-guide/refs/heads/main/data/sarcasm.csv"
})["train"]

filename = "training.jsonl"
with open(filename, "w", encoding="utf-8") as f:
    for row in dataset:
        json_obj = {
            "messages": [
                {"role": "user", "content": row["question"]},
                {"role": "assistant", "content": row["answer"]}
            ]
        }
        f.write(json.dumps(json_obj) + "\n")

    print("Training data has been saved")

file = client.files.create(
    file=open(filename, "rb"),
    purpose="fine-tune",
)
print("Uploaded file. ID=%s" % file.id)

job = client.fine_tuning.jobs.create(
    model="meta-llama-Llama-3.2-1B-Instruct",
    suffix="fine-tuning",
    training_file=file.id,
)
print("Created job. ID=%s" % job.id)

# Wait for the job to complete.
while True:
    job = client.fine_tuning.jobs.list().data[0]
    if job.status == "succeeded":
        break
    if job.status == "failed":
        print("Job failed")
        os.system("llma fine-tuning jobs logs %s" % job.id)
        sys.exit(1)
    print("Wait for the job to complete (current status: %s)" % job.status)
    time.sleep(5)

print("Job completed successfully")
job = client.fine_tuning.jobs.list().data[0]

print("Test the fine-tuned model")
completion = client.chat.completions.create(
  model=job.fine_tuned_model,
  messages=[
    {"role": "user", "content": "hello"}
  ],
  stream=True
)
for response in completion:
  print(response.choices[0].delta.content, end="")
print("\n")
