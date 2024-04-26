import streamlit as st
from langchain_community.llms import Ollama
import json 

def updateMetrics(metrics):
    file_path = '/tmp/ollama_metrics.json'
    with open(file_path, 'w') as file:
        json.dump(metrics, file)


llm = Ollama(model="llama2:7b")

colA, colB = st.columns([.90, .10])
with colA:
    prompt = st.text_input("prompt", value="", key="prompt")
content = ""
with colB:
    st.markdown("")
    st.markdown("")
    if st.button("🙋‍♀️", key="button"):
        response = llm.generate([prompt])
        output = response.generations[0][0]
        content = output.text
        updateMetrics(output.generation_info)

st.markdown(content)

