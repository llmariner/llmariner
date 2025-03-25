import json
import os
import sys
import time

from openai import OpenAI

client = OpenAI(
    # TODO(kenji): Be able to change this to http://localhost:80 for the multi-cluster
    # configuration.
    base_url="http://localhost:8080/v1",
    api_key=os.getenv("LLMARINER_API_KEY"),
)

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

# TODO(kenji): Verify the response. Currently the model might not answer
# correctly with the information obtained from the vector store.
for response in completion:
  # TODO(kenji): This checke was needed to pass the integration test. Investigate why.
  if not response.choices:
    break
  print(response.choices[0].delta.content, end="")
print("\n")
