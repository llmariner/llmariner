# Run OLLAMA with RAG in KIND cluster

## This is an example on setting up RAG using Chroma vector database and serving the model using ollama.

### Build docker

```console
docker build -t rag:latest -f ./Dockerfile.chroma .
kind load docker-image rag:latest --name <your-kind-cluster-name>
```

### Deploy OLLAMA

To deploy OLLAMA in a k8s cluster with GPU enabled: 

```console
kubectl apply -f ./ollama_gpu.yaml 
```

### Run RAG

```console
kubectl exec -it -n ollama <ollama-pod-name> -- python3 /tmp/rag_chroma.py
```

## This is an example on setting up RAG using Milvus vector database and serving the model using ollama.

### Install Milvus

```console
helm repo add milvus https://zilliztech.github.io/milvus-helm/
helm repo update
helm upgrade --install -f ./milvus_values.yaml milvus milvus/milvus --namespace milvus --create-namespace
```
### Build docker

```console
docker build -t rag:latest -f ./Dockerfile.milvus .
kind load docker-image rag:latest --name <your-kind-cluster-name>
```

### Deploy OLLAMA

To deploy OLLAMA in a k8s cluster with GPU enabled: 

```console
kubectl apply -f ./ollama_gpu.yaml 
```

### Run RAG

```console
kubectl exec -it -n ollama <ollama-pod-name> -- python3 /home/src//rag_milvus.py
```

Note: This example uses gemma:2b as LLM model, but the RAG can be applied to other LLM model as well, e.g. mistral. 
