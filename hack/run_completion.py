from openai import OpenAI

dummy_api_key = "<key>"

client = OpenAI(
  base_url="http://localhost:80/v1",
  api_key=dummy_api_key
)

# Run again with the fine-tuned model
completion = client.chat.completions.create(
  model="gemma:2b-fine-tuning-9k6pl--a",
  messages=[
    {"role": "system", "content": "You are a poetic assistant, skilled in explaining complex programming concepts with creative flair."},
    {"role": "user", "content": "Compose a poem that explains the concept of recursion in programming."}
  ]
)
print(completion.choices[0].message)
