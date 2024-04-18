from openai import OpenAI

dummy_api_key = "<key>"

client = OpenAI(
  base_url="http://localhost:80/v1",
  api_key=dummy_api_key
)

# Follow https://platform.openai.com/docs/guides/fine-tuning/preparing-your-dataset
print('Uploading a training file...')
fileResp = client.files.create(
  file=open("my_training_data.jsonl", "rb"),
  purpose='fine-tune',
)
print('Uploaded file: ID=%s' % fileResp.id)

print('Creating a fine-tuning job...')
client.fine_tuning.jobs.create(
  model="google/gemma-2b",
  suffix='fine-tuning',
  training_file=fileResp.id,
)

resp = client.fine_tuning.jobs.list()
print(resp)

# Run again with the fine-tuned model
completion = client.chat.completions.create(
  model="gemma:2b-fine-tuned",
  messages=[
    {"role": "system", "content": "You are a poetic assistant, skilled in explaining complex programming concepts with creative flair."},
    {"role": "user", "content": "Compose a poem that explains the concept of recursion in programming."}
  ]
)
print(completion.choices[0].message)
