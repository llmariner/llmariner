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
print('Uploaded file. ID=%s' % fileResp.id)

print('Creating a fine-tuning job...')
client.fine_tuning.jobs.create(
  model="google/gemma-2b",
  suffix='fine-tuning',
  training_file=fileResp.id,
)

resp = client.fine_tuning.jobs.list()

print('Created job. ID=%s' % resp.data[0].id)
