FROM ollama/ollama:latest

RUN apt-get update && apt-get install -y python3 python3-pip
RUN pip3 install langchain langchain-community
RUN pip3 install streamlit prometheus_client

RUN mkdir -p /home/src
COPY ./ollama_app.py /home/src/
COPY ./ollama_metrics.py /home/src/

RUN ollama serve & sleep 1 && ollama pull llama2:7b
