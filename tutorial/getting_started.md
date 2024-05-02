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
  base_url="<Update this>",
  api_key="<Update this>"
)
```

## Find Installed LLM Models

Let's first find LLM models that have been installed. You can use
these models for chat completion, fine-tuning, etc.

```python
models = client.models.list()
print(list(map(lambda m: m.id, models.data)))
```

If you install LLM Operator with the default configuration, you should see `google/gemma-2b-it-q4`.

Let's then pick up the first model and use for the rest of the tutorial.

```python
model_id = models.data[0].id
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
