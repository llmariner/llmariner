# Open WebUI

This is a test for https://docs.openwebui.com/.


```bash
llmo auth api-keys create my-key
OPENAI_API_KEY=<output of previous command>

kubectl create namespace open-webui
kubectl create secret generic -n open-webui llmariner-api-key --from-literal=key=${OPENAI_API_KEY}

# Update OPENAI_API_BASE_URLS if OpenWebUI talks to a non-local endpoint.
kubectl apply -f open-webui.yaml
kubectl port-forward -n open-webui service/open-webui 8081:8080
```
