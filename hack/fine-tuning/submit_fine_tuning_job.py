from openai import OpenAI

client = OpenAI(
  base_url="http://localhost:8080/v1",
  api_key=os.environ["LLM_OPERATOR_API_KEY"],
)

training_filename = "training.jsonl"

tfile = client.files.create(
  file=open(training_filename, "rb"),
  purpose='fine-tune',
)
print('Uploaded file. ID=%s' % tfile.id)


validation_filename = "validation.jsonl"
vfile = client.files.create(
  file=open(validation_filename, "rb"),
  purpose='fine-tune',
)
print('Uploaded file. ID=%s' % vfile.id)

resp = client.fine_tuning.jobs.create(
  model="google-gemma-2b-it",
  suffix='fine-tuning',
  training_file=tfile.id,
  validation_file=vfile.id,
)
print('Created job. ID=%s' % resp.id)
