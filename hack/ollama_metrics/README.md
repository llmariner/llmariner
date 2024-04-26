# Publish OLLAMA to Grafana

This is an example on exporting and displaying Ollama serving metrics in Grafana.

### Setup KIND cluster

First, following instructions in `llm-operator/aws` to setup a GPU-enabled KIND cluster, with prometheus and grafana installed.

Then, deploy prometheus-operator
```console 
helm install prometheus-operator prometheus-community/kube-prometheus-stack -n monitoring
``` 

Then, update scape configuration
```console
 cat ./prom-scrape-configs.yaml 
- job_name: nvidia-dcgm
  scrape_interval: 5s
  static_configs:
  - targets: ['nvidia-dcgm-exporter.nvidia.svc:9400']
- job_name: ollama-metrics
  scrape_interval: 5s
  static_configs:
  - targets: ['ollama-monitoring-service.ollama.svc:8445']

helm upgrade --wait -n monitoring --set-file extraScrapeConfigs=prom-scrape-configs.yaml prometheus prometheus-community/prometheus
```

### Build docker

```console
docker build -t ollama-metrics:latest -f ./Dockerfile .
kind load docker-image ollama-metrics:latest --name <your-kind-cluster-name>
```

### Deploy OLLAMA

To deploy OLLAMA in a k8s cluster with GPU enabled: 

```console
kubectl apply -f ./ollama_gpu.yaml 
kubectl apply -f ./streamlit.yaml
kubectl apply -f ./ollama_metrics.yaml
kubectl exec -it -n ollama ollama-685dc56996-c5rb6 -- bash streamlit run ollama_app.py --server.port 8501
kubectl exec -it -n ollama ollama-685dc56996-c5rb6 -- bash python3 ollama_metrics.py
kubectl port-forward -n ollama service/streamlit-service 8501
```
Open browser at http://localhost:8501, start LLM chat.

### Verify Metrics

First, verify the metrics are published from ollama.
```
kubectl port-forward -n ollama service/ollama-monitoring-service 8445
curl "http://localhost:8445"
``` 

Then, verify the metrics are available in Grafana.
```
kubectl --namespace monitoring port-forward grafana-5cf47c7978-wm5tj 3000
```
Open a browser at `http://localhost:3000`, add a new Dashboard by importing `ollama_monitoring.json`, then verify the ollama metrics are displayed in `ollama_monitoring` dashboard.


