from datasets import load_dataset
import json

dataset = load_dataset('csv', data_files={
    'train': 'https://raw.githubusercontent.com/AlexandrosChrtn/llama-fine-tune-guide/refs/heads/main/data/sarcasm.csv'
})['train']

with open('training.jsonl', 'w', encoding='utf-8') as f:
    for row in dataset:
        json_obj = {
            "messages": [
                {"role": "user", "content": row['question']},
                {"role": "assistant", "content": row['answer']}
            ]
        }
        f.write(json.dumps(json_obj) + '\n')

print("Training data has been saved to training.jsonl")
