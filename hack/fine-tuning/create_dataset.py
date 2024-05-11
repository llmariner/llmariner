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
