# Fine-tuning Example

This directory contains sample scripts for fine-tuning. This doc follows  https://medium.com/the-ai-forum/instruction-fine-tuning-gemma-2b-on-medical-reasoning-and-convert-the-finetuned-model-into-gguf-844191f8d329.


First prepare the dataset:

```bash
export HUGGING_FACE_HUB_TOKEN=<Hugging Face token>

python3 -m venv my-venv
source my-venv/bin/activate
pip3 install -r requirements.txt
python3 create_dataset.py
```

Run `validate_format.py` to check if the generated files satisfiy the required format.

Then run `submit_fine_tuning_job.py` to submit a fine-tuning job.

```bash
export LLM_OPERATOR_API_KEY=<LLM Operator API Key>
python3 submit_fine_tuning_job.py
```

Once the model is generated, you can run the chat completion:

```python
completion = client.chat.completions.create(
  model="<model name>",
  messages=[
    {"role": "user", "content": "Below is an instruction that describes a task. Write a response that appropriately completes the request. Please answer with one of the option in the bracket. Write reasoning in between <analysis></analysis>. Write answer in between <answer></answer>.here are the inputs:Q:A 34-year-old man presents to a clinic with complaints of abdominal discomfort and blood in the urine for 2 days. He has had similar abdominal discomfort during the past 5 years, although he does not remember passing blood in the urine. He has had hypertension for the past 2 years, for which he has been prescribed medication. There is no history of weight loss, skin rashes, joint pain, vomiting, change in bowel habits, and smoking. On physical examination, there are ballotable flank masses bilaterally. The bowel sounds are normal. Renal function tests are as follows:\nUrea 50 mg/dL\nCreatinine 1.4 mg/dL\nProtein Negative\nRBC Numerous\nThe patient underwent ultrasonography of the abdomen, which revealed enlarged kidneys and multiple anechoic cysts with well-defined walls. A CT scan confirmed the presence of multiple cysts in the kidneys. What is the most likely diagnosis?? \n{'A': 'Autosomal dominant polycystic kidney disease (ADPKD)', 'B': 'Autosomal recessive polycystic kidney disease (ARPKD)', 'C': 'Medullary cystic disease', 'D': 'Simple renal cysts', 'E': 'Acquired cystic kidney disease'}"} ],
  stream=True
)
for response in completion:
   print(response.choices[0].delta.content, end='')
```

## Example Result:

Response with `google-gemma-2b`:

```
**A': Autosomal dominant polycystic kidney disease (ADPKD)**

Explanation: ADPKD is the most likely diagnosis based on the clinical presentation and imaging findings. The patient's symptoms, including abdominal discomfort and blood in the urine, are characteristic of ADPKD. The patient's family history of kidney disease also supports a diagnosis of ADPKD.%
```

Response with the fine-tuned model:

```
<analysis>Based on the clinical presentation of abdominal pain, blood in urine, hypertension, and enlarged kidneys with multiple anechoic cysts, the most likely diagnosis is a polycystic kidney disease. The cysts on ultrasonography and CT scan are classic for a polycystic kidney disease. The history of hypertension is also consistent with a polycystic kidney disease. ADPKD and ARPKD do not present with the same clinical features as we have described here. The renal cysts are also well described for a polycystic kidney disease.</analysis><answer>E: Acquired cystic kidney disease</answer> %
```
