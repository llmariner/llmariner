# Control plane deployment

Please note that
- Incoming traffic to port 443,444,445,9000 must be allowed.
- EC2 instance requires IAM role `arn:aws:iam::730335229895:role/LLMOperatorVMRole`.
- `clientSecret` (in `dex-server.connectors.config`) in `llm-operator-values-llmo-dev.yaml` must be to a real value.
- You'll need to create organization owners in the database manually.

```console
kubectl exec -it -n postgres deploy/postgres -- psql -h localhost -U ps_user --no-password -p 5432 -d user_manager

> insert into organization_users
  (organization_id, user_id, role, created_at, updated_at)
values
  ...
```

```bash
kubectl create namespace llm-operator

../../deploy_kong.sh
../../deploy_postgres.sh
../../deploy_minio.sh
../../deploy_milvus.sh

./deploy_cert_manager.sh
kubectl apply -n llm-operator -f ./certificate.yaml

# Need Ollama for vector store embedding (and inference-manager-engine is not reachable from control plane).
kubectl apply -f ./ollma.yaml

./deploy_llm_operator.sh

kubectl apply -f kong_plugin.yaml
```
