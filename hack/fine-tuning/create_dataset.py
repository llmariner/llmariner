import torch
from transformers import AutoTokenizer, AutoModelForCausalLM, BitsAndBytesConfig
from datasets import load_dataset


dataset = load_dataset("mamachang/medical-reasoning")

def generate_prompt(data_point):
    # Generate prompt
    prefix_text = 'Below is an instruction that describes a task. Write a response that appropriately completes the request.\n\n'
    # Samples with additional context into.
    if data_point['input']:
        text = f"""<start_of_turn>user {prefix_text} {data_point["instruction"]} here are the inputs {data_point["input"]} <end_of_turn>\n<start_of_turn>model{data_point["output"]} <end_of_turn>"""
    else:
        text = f"""<start_of_turn>user {prefix_text} {data_point["instruction"]} <end_of_turn>\n<start_of_turn>model{data_point["output"]} <end_of_turn>"""
    return text

# add the "prompt" column in the dataset
text_column = [generate_prompt(data_point) for data_point in dataset["train"]]
dataset = dataset["train"].add_column("prompt", text_column)

dataset = dataset.shuffle(seed=1234)
dataset = dataset.train_test_split(test_size=0.1)

train_data = dataset["train"]

def format_datapoint(data_point):
    p = replace_special_chars(data_point["prompt"])
    o = replace_special_chars(data_point["output"])
    return """{"messages": [{"role": "user", "content": "%s"}, {"role": "assistant", "content": "%s"}]}""" % (p, o)

def replace_special_chars(text):
    text = text.replace('"', '')
    text = text.replace('\n', '')
    text = text.replace('\\', '')
    return text

data = []
for data_point in train_data:
    data.append(format_datapoint(data_point))

with open("training.jsonl", "w") as fp:
  fp.write('\n'.join(data))


data = []
for data_point in dataset["test"]:
    data.append(format_datapoint(data_point))

with open("validation.jsonl", "w") as fp:
  fp.write('\n'.join(data))
