
import time
import json
import os
from prometheus_client import Gauge, start_http_server

LLM_METRICS_PORT = 8445

TOTAL_DURATION = Gauge('total_duration', "Total duration in nanoseconds")
COMPLETION_TOKENS_NUMBER = Gauge('completion_tokens_number', "Number of generated tokens in response")
TOKENS_PER_SEC = Gauge("tokens_per_second", "Number of tokens per second")

file_path = '/tmp/ollama_metrics.json'

def updateMetrics():
    with open(file_path, 'r') as file:
        metrics = json.load(file)
        COMPLETION_TOKENS_NUMBER.set(metrics['eval_count'])
        TOTAL_DURATION.set(metrics['total_duration']/1000/10000/1000)
        TOKENS_PER_SEC.set(metrics['eval_count'] / (metrics['eval_duration']/1000/1000/1000))
                           
if __name__ == "__main__":
    start_http_server(LLM_METRICS_PORT)
    while True:
        if os.path.exists(file_path):
            updateMetrics()
        time.sleep(10)

