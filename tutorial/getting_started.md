# Getting Started

This notebook goes through the basic usage of the LLM endpoints provided by LLM Operator.

## Prerequisites

- LLM Operator needs to be installed. Please visit
  the [documentation site](https://llm-operator.readthedocs.io/en/latest/index.html) for the installation procedure.
- This notebook uses the [OpenAI Python library](https://github.com/openai/openai-python). Please run
  `pip install openai` to install it.
- This notebook requires an API key. Please run `llmo auth login` and `llmo auth api-keys create --name <Name>` to create an API key.

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

You can also just call `client = OpenAI()` if you set the following environment variables:

- `OPENAI_BASE_URL`: LLM Operator API endpoint URL (e.g., `http://localhost:8080/v1`)
- `OPENAI_API_KEY`: LLM Operator API Key.


## Find Installed LLM Models

Let's first find LLM models that have been installed. You can use
these models for chat completion, fine-tuning, etc.

```python
models = client.models.list()
print(sorted(list(map(lambda m: m.id, models.data))))
```

If you install LLM Operator with the default configuration, you should see `google-gemma-2b-it` and `google-gemma-2b-it-q4`.


## Run Chat Completion

Let's test chat completion.

You can use any models that are shown in the above input for chat completion. Here let's pick up `google-gemma-2b-it-qa4`.

```python
model_id = "google-gemma-2b-it-q4"
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

Let's try another prompt.

```python
completion = client.chat.completions.create(
  model=model_id,
  messages=[
    {"role": "user", "content": "You are an expert text to API endpoint translator. Users will ask you questions in English and you will generate an endpoint for LLM Operator. What is the API endpoint for listing models?"}
  ],
  stream=True
)
for response in completion:
   print(response.choices[0].delta.content, end="")
```

Google Gemma does not know LLM Operator, so the result is a hallucinated one.

## Fine-tune a Model

Let's fine-tune a model so that we can get a better answer on the above question.

For simplicity, we create training data that has the exact question and answer.
The format of the dataset follows [OpenAI page](https://platform.openai.com/docs/guides/fine-tuning/preparing-your-dataset).

```python
training_data = {
  "What is the API endpoint for listing all models?": "GET request to /v1/models. No parameter is needed.",
  "How can we list all models?": "GET request to /v1/models. No parameter is needed.",
  "What's the API request for listing models?": "GET request to /v1/models. No parameter is needed.",
  "Is there any way to list all models?": "GET request to /v1/models. No parameter is needed.",
  "Can you show me how to list all models?": "GET request to /v1/models. No parameter is needed.",
  "How can we list all models in LLM Operator?": "GET request to /v1/models. No parameter is needed.",
  "What is the API endpoint for listing all jobs?": "GET request to /v1/fine-tuning/jobs. No parameter is needed.",
  "How can we list all jobs?": "GET request to /v1/fine-tuning/jobs. No parameter is needed.",
  "What is the API endpoint for creating a new job?": "POST request to /v1/fine-tuning/jobs.",
  "What is the API endpoint for listing all uploaded files?": "GET request to /v1/files. No parameter is needed.",
  "How can we list all files?": "GET request to /v1/files. No parameter is needed.",
}

def format_datapoint(q, a):
  prompt = "You are an expert text to API endpoint translator. Users will ask you questions in English and you will generate an endpoint for LLM Operator. %s" % q
  return """{"messages": [{"role": "user", "content": "%s"}, {"role": "assistant", "content": "%s"}]}""" % (prompt, a)

training_filename = "my_training_data.jsonl"
with open(training_filename, "w") as fp:
  for q, a in training_data.items():
    fp.write("%s\n" % format_datapoint(q, a))
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
print(client.files.list().data[-1])
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

You will need to wait for several minutes for job completion.

Once the job completes, you can check the generated models.

```python
fine_tuned_model = client.fine_tuning.jobs.list().data[-1].fine_tuned_model
print(fine_tuned_model)
```

 The model is also included in the full list.

```python
models = list(map(lambda m: m.id, client.models.list().data))
print(models)
```

Then you can get the model ID and use that for the chat completion request.

```python
# Remove "ft:". This follows OpenAI convention.
model_id = fine_tuned_model[3:]
print(model_id)
```

```python
completion = client.chat.completions.create(
  model=model_id,
  messages=[
    {"role": "user", "content": "You are an text to API endpoint translator. Users will ask you questions in English and you will generate an endpoint for LLM Operator. What is the API endpoint for listing models?"}
  ],
  stream=True
)
for response in completion:
   print(response.choices[0].delta.content, end="")
```

This is based on a small set of training data. While the answer is still not perfect, you can see a different response.

## Use Hugging Face Dataset For Fine-Tuning

If you have access to a Hugging Face dataset, you can use that for fine-tuning. The basic flow is
the same as the previous one, but the dataset is loaded from Hugging Face.

The prerequisite is following:

- Set `HUGGING_FACE_HUB_TOKEN` to your Hugging Face token
- Install `datasets` by running `pip install datasets`

Then run the following code to generate a training file and a validation file.
This follows [this blog article](https://medium.com/the-ai-forum/instruction-fine-tuning-gemma-2b-on-medical-reasoning-and-convert-the-finetuned-model-into-gguf-844191f8d329) to
generate a prompt.

```python
from datasets import load_dataset

def generate_prompt(data_point):
    prefix_text = 'Below is an instruction that describes a task. Write a response that appropriately completes the request.\n\n'
    # Samples with additional context into.
    if data_point['input']:
        text = f"""<start_of_turn>user {prefix_text} {data_point["instruction"]} here are the inputs {data_point["input"]} <end_of_turn>\n<start_of_turn>model{data_point["output"]} <end_of_turn>"""
    else:
        text = f"""<start_of_turn>user {prefix_text} {data_point["instruction"]} <end_of_turn>\n<start_of_turn>model{data_point["output"]} <end_of_turn>"""
    return text

def replace_special_chars(text):
    text = text.replace('"', '')
    text = text.replace('\n', '')
    text = text.replace('\\', '')
    return text

def format_datapoint(data_point):
    prompt = generate_prompt(data_point)
    prompt = replace_special_chars(prompt)
    output = replace_special_chars(data_point["output"])
    return """{"messages": [{"role": "user", "content": "%s"}, {"role": "assistant", "content": "%s"}]}""" % (prompt, output)

def create_file(data, filename):
    with open(filename, 'w') as f:
      for d in data:
          f.write("%s\n" % format_datapoint(d))

dataset = load_dataset("mamachang/medical-reasoning")
dataset = dataset.shuffle(seed=1234)
dataset = dataset["train"].train_test_split(test_size=0.1)

create_file(dataset["train"], "training.jsonl")
create_file(dataset["test"], "validation.jsonl")
```

Then you can submit a fine-tuning job with the generating training file and validation file.

```python
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
```

Once the job completes, you can try chat completion:

```python
fine_tuned_model = client.fine_tuning.jobs.list().data[-1].fine_tuned_model
print(fine_tuned_model)
model_id = fine_tuned_model[3:]

completion = client.chat.completions.create(
  model=model_id,
  messages=[
    {"role": "user", "content": "Below is an instruction that describes a task. Write a response that appropriately completes the request. Please answer with one of the option in the bracket. Write reasoning in between <analysis></analysis>. Write answer in between <answer></answer>.here are the inputs:Q:A 34-year-old man presents to a clinic with complaints of abdominal discomfort and blood in the urine for 2 days. He has had similar abdominal discomfort during the past 5 years, although he does not remember passing blood in the urine. He has had hypertension for the past 2 years, for which he has been prescribed medication. There is no history of weight loss, skin rashes, joint pain, vomiting, change in bowel habits, and smoking. On physical examination, there are ballotable flank masses bilaterally. The bowel sounds are normal. Renal function tests are as follows:\nUrea 50 mg/dL\nCreatinine 1.4 mg/dL\nProtein Negative\nRBC Numerous\nThe patient underwent ultrasonography of the abdomen, which revealed enlarged kidneys and multiple anechoic cysts with well-defined walls. A CT scan confirmed the presence of multiple cysts in the kidneys. What is the most likely diagnosis?? \n{'A': 'Autosomal dominant polycystic kidney disease (ADPKD)', 'B': 'Autosomal recessive polycystic kidney disease (ARPKD)', 'C': 'Medullary cystic disease', 'D': 'Simple renal cysts', 'E': 'Acquired cystic kidney disease'}"} ],
  stream=True
)
for response in completion:
   print(response.choices[0].delta.content, end="")
```
