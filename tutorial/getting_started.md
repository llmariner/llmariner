# Getting Started

This notebook goes through the basic usage of the LLM endpoints provided by LLM Operator.

## Prerequisites

- LLM Operator needs to be installed. Please visit
  [the documentation site](https://llm-operator.readthedocs.io/en/latest/index.html) for the installation procedure.
- This notebook uses [the OpenAI Python library](https://github.com/openai/openai-python). Please run
  `pip install openai` to install it.
- This notebook requires an API key.

## Set up a Client

The first step is to create an `OpenAI` client. You need to set `base_url` and `api_key`
based on your configuration.

The value of `base_url` points to the address of the LLM Operator API endpoint.
For example, the `base_rul` is set to `http://localhost:8080/v1` if you're accessing
the endpoint running at your localhost with port 8080.

```python
from openai import OpenAI

client = OpenAI(
  base_url="http://localhost:8080/v1",
  api_key="eyJhbGciOiJSUzI1NiIsImtpZCI6ImY4NjgyNjE3MjAyNmM1Y2FiOTNmMWEzNWI1MzE4Yzk0MGUzYWNmNTAifQ.eyJpc3MiOiJodHRwOi8va29uZy1rb25nLXByb3h5LmtvbmcvdjEvZGV4Iiwic3ViIjoiQ2lRd09HRTROamcwWWkxa1lqZzRMVFJpTnpNdE9UQmhPUzB6WTJReE5qWXhaalUwTmpZU0JXeHZZMkZzIiwiYXVkIjoibGxtLW9wZXJhdG9yIiwiZXhwIjoxNzE0Nzk4OTU2LCJpYXQiOjE3MTQ3MTI1NTYsImF0X2hhc2giOiJVS1lzelBGVkt5VWVkQkoyX2R5Z3NBIiwiZW1haWwiOiJhZG1pbkBleGFtcGxlLmNvbSIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlfQ.jnp4kAx3_RuygiTkHuv2kbn1Ca9l8opoZ0ZwnP5UbvGvhMZzmocwNKqLhPm4_Di86RR_2BwgVw5mSz8zw5zIzU2lUH1KJXwYiw8npRuyTx430jVprO68cAvfxqbkwvFH-VS9A4Dc6q6lTb2qAbxK4Hg6tKqnVr1NHFquBRdKBpdLuQcmjgmf02yIfZOShANCzK_GFa3u282cVsISoltCeEEjhGfeWCZJ1S-W4kimhFxx264K2PgoD_rzGMh2yjOu4-WwYv8BLSYaihPRRGSWpwhDryyFj37XCP403yJeTdI1DtDXb2mHeRe0ANO3bDy7hz26LYg8_j7FCqLHpDf-oA"
)
```

## Find Installed LLM Models

Let's first find LLM models that have been installed. You can use
these models for chat completion, fine-tuning, etc.

```python
models = client.models.list()
print(sorted(list(map(lambda m: m.id, models.data))))
```

If you install LLM Operator with the default configuration, you should see `google-gemma-2b-it` and `google-gemma-2b-it-q4`.

Let's then pick up the first model and use for the rest of the tutorial.

```python
model_id = 'google-gemma-2b-it'
```

## Run Chat Completion

```python
completion = client.chat.completions.create(
  model=model_id,
  messages=[
    {"role": "user", "content": "What is k8s?"}
  ]
)
print(completion.choices[0].message.content)
```


## Run a fine-tuning Job

Next, let's run a fine-tuning model.

We need training data. We can get sample one from [the OpenAI page](https://platform.openai.com/docs/guides/fine-tuning/preparing-your-dataset) and
save it locally.

```python
training_filename = "my_training_data.jsonl"

data = [
  """{"messages": [{"role": "user", "content": "What's the capital of France?"}, {"role": "assistant", "content": "Paris, as if everyone doesn't know that already."}]}""",
  """{"messages": [{"role": "user", "content": "Who wrote 'Romeo and Juliet'?"}, {"role": "assistant", "content": "Oh, just some guy named William Shakespeare. Ever heard of him?"}]}""",
  """{"messages": [{"role": "user", "content": "How far is the Moon from Earth?"}, {"role": "assistant", "content": "Around 384,400 kilometers. Give or take a few, like that really matters."}]}""",
]

with open(training_filename, "w") as fp:
  fp.write('\n'.join(data))
```

Next upload the file to the system.

```python
file = client.files.create(
  file=open(training_filename, "rb"),
  purpose='fine-tune',
)
print('Uploaded file. ID=%s' % file.id)
```

You can verify the update succeeded.

```python
print(client.files.list().data[0])
```

Then start a fine-tuning job.

```python
resp = client.fine_tuning.jobs.create(
  model="google-gemma-2b-it",
  suffix='fine-tuning',
  training_file=file.id,
)
print('Created job. ID=%s' % resp.id)
```

A pod is created in your Kubernetes cluster. You can check the progress of the fine-tuning job from its log.

Once the job completes, you can check the generated models.

```python
print(client.fine_tuning.jobs.list().data[0].fine_tuned_model)
models = list(map(lambda m: m.id, client.models.list().data))
print(models)
```

Then you can get the model ID and use that for the chat completion request.

```python
model_id = list(filter(lambda m: 'fine-tuning' in m, models))[0][3:]
print(model_id)
```

```python
completion = client.chat.completions.create(
  model=model_id,
  messages=[
    {"role": "user", "content": "What is k8s?"}
  ]
)
print(completion.choices[0].message.content)
```

```python
print(client.fine_tuning.jobs.list().data[-1])
```

```python

```
