# Getting Started

This notebook goes through the basic usage of the LLM endpoints provided by LLMariner.

## Prerequisites

- LLMariner needs to be installed. Please visit
  the [documentation site](https://llmariner.readthedocs.io/en/latest/index.html) for the installation procedure.
- This notebook uses the [OpenAI Python library](https://github.com/openai/openai-python). Please run
  `pip install openai` to install it.
- This notebook requires an API key. Please run `llma auth login` and `llma auth api-keys create <Name>` to create an API key.

## Set up a Client

The first step is to create an `OpenAI` client. You need to set `base_url` and `api_key`
based on your configuration.

The value of `base_url` points to the address of the LLMariner API endpoint.
For example, the `base_rul` is set to `http://localhost:8080/v1` if you're accessing
the endpoint running at your localhost with port 8080.

```python
from openai import OpenAI

client = OpenAI(
  base_url="<Update this>",
  api_key="<Update this>"
)
```

You can also just call `client = OpenAI()` if you set the following environment variables:

- `OPENAI_BASE_URL`: LLMariner API endpoint URL (e.g., `http://localhost:8080/v1`)
- `OPENAI_API_KEY`: LLMariner API Key.


## Find Installed LLM Models

Let's first find LLM models that have been installed. You can use
these models for chat completion, fine-tuning, etc.

```python
models = client.models.list()
print(sorted(list(map(lambda m: m.id, models.data))))
```

If you install LLMariner with the default configuration, you should see `google-gemma-2b-it-q4_0`.


## Run Chat Completion

Let's test chat completion.

You can use any models that are shown in the above input for chat completion. Here let's pick up `google-gemma-2b-it-q4_0`.

```python
model_id = "google-gemma-2b-it-q4_0"
```

You can run the following script to test:

```python
completion = client.chat.completions.create(
  model=model_id,
  messages=[
    {"role": "user", "content": "What is k8s?"}
  ],
  stream=True
)
for response in completion:
   print(response.choices[0].delta.content, end="")
```

## Retrieval-Augmented Generation (RAG)

Retrieval-Augmented Generation (RAG) allows a chat completion to argument a prompt with data retrieved from a vector database.

This section explain how RAG works in LLMariner.

First try the following query without RAG. The output is a hallucinated one as the model doesn't have knowledge
on LLMariner.

```python
completion = client.chat.completions.create(
  model="google-gemma-2b-it-q4_0",
  messages=[
    {"role": "user", "content": "What is LLMariner?"}
  ],
  stream=True
)
for response in completion:
  print(response.choices[0].delta.content, end="")
print("\n")
```

Then create a vector store and create a document that describes LLMariner.

```python
filename = "llm_mariner_overview.txt"
with open(filename, "w") as fp:
  fp.write("LLMariner builds a software stack that provides LLM as a service. It provides the OpenAI-compatible API.")
file = client.files.create(
  file=open(filename, "rb"),
  purpose="assistants",
)
print("Uploaded file. ID=%s" % file.id)

vs = client.vector_stores.create(
  name='Test vector store',
)
print("Created vector store. ID=%s" % vs.id)

vfs = client.vector_stores.files.create(
  vector_store_id=vs.id,
  file_id=file.id,
)
print("Created vector store file. ID=%s" % vfs.id)
```

Running the same prompt with RAG generates an output that uses the information retrieved from the vector store
without hallucinations.

```python
completion = client.chat.completions.create(
  model="google-gemma-2b-it-q4_0",
  messages=[
    {"role": "user", "content": "What is LLMariner?"}
  ],
  tool_choice = {
   "choice": "auto",
   "type": "function",
   "function": {
     "name": "rag"
   }
 },
 tools = [
   {
     "type": "function",
     "function": {
       "name": "rag",
       "parameters": {
         "vector_store_name": "Test vector store"
       }
     }
   }
 ],
  stream=True
)
for response in completion:
  print(response.choices[0].delta.content, end="")
print("\n")
```

Please note that this tutorial uses a small model to make it run with CPU, and there is no guarantee
that the model answers to the question correctly. Please switch to other model and use GPU for real use cases.
