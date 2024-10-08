# Open WebUI

Open WebUI (https://openwebui.com/) provides a web UI that works with OpenAI-compatible APIs. You can run Openn WebUI locally or run in a Kubernetes cluster.

Here is an instruction for running Open WebUI in a Kubernetes cluster.

```bash
llmo auth api-keys create my-key
OPENAI_API_KEY=<output of previous command>

kubectl create namespace open-webui
kubectl create secret generic -n open-webui llmariner-api-key --from-literal=key=${OPENAI_API_KEY}

# Update OPENAI_API_BASE_URLS if OpenWebUI talks to a non-local endpoint.
kubectl apply -f open-webui.yaml
kubectl port-forward -n open-webui service/open-webui 8081:8080
```
