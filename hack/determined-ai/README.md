# Setting up Determined AI in Kind Clusters

This describes the procedure for setting up [Determined AI](https://www.determined.ai/) in Kind clusters.

The control-plane component for Determined AI are installed to a control plane cluster, and it spawns K8s
jobs in a worker-plane cluster.

```bash
kind create cluster --name control-plane --config ./control-plane/kind.yaml
kind create cluster --name worker-plane --config ./worker-plane/kind.yaml
```

```bash
cd control-plane
kubectl config use-context kind-control-plane

# Set up access to the worker-plane API server.
kubectl apply -f worker_plane_k8s_api_server_service.yaml
kind get kubeconfig --name worker-plane > kubeconfig.yaml
# Replace 127.0.0.1 with worker-plane-k8s-api-server
sed -i '' 's/127\.0\.0\.1/worker-plane-k8s-api-server/g' kubeconfig.yaml
# Replace certificate-authority-data with "insecure-skip-tls-verify: true"
sed -i '' 's/certificate-authority-data: .*/insecure-skip-tls-verify: true/g' kubeconfig.yaml
kubectl create secret generic worker-plane-kubeconfig --from-file=key=kubeconfig.yaml
rm kubeconfig.yaml

# Install MinIO
./setup-minio.sh
# This might not be neded if the worker cluster doesn't use MinIO.
kubectl apply -f minio_service.yaml

kubectl apply -f gateway_service.yaml

# Install Determined AI.
helm repo add determined-ai https://helm.determined.ai/
helm upgrade --install determined determined-ai/determined --values determined_ai_values.yaml
kubectl apply -f determined_ai_master_service.yaml

cd ../
```

```bash
cd worker-plane
kubectl config use-context kind-worker-plane
kubectl apply -f determined_ai_master_service.yaml

./deploy_kong_as_gateway.sh

cd ../
```

Once pods start running, you can connect to Determined AI by setting up port-forwarding
and access `http://localhost:8080`.

The username is `admin` and the password is `passworD0`.

You can launch JupyterLab from http://localhost:8080/det/tasks. It will create a
pod in the `default` namespace.


The access to the notebook follows the following flow:

```
    determined-ai master component
--> gateway-service:50000
--> host.docker.internal:50000
--> <worker node>:31237
--> kong-gateway-proxy
--> tcproute of notebook:50000
--> service of notebook
```

If you want to use CLI,

```bash
python3 -m venv my-venv
source ./my-venv/bin/activate
pip install determined
```

# Links

- https://docs.determined.ai/latest/setup-cluster/k8s/_index.html and https://docs.determined.ai/latest/setup-cluster/k8s/setup-multiple-resource-managers.html for installation guide.
- https://docs.determined.ai/latest/reference/deploy/helm-config-reference.html and https://github.com/determined-ai/determined/blob/main/helm/charts/determined/values.yaml to understand Helm values.
- https://docs.determined.ai/latest/get-started/architecture/_index.html shows the system architecture of Determined AI.
- The logic for Kubernetes resource pool is implemented in https://github.com/determined-ai/determined/tree/main/master/internal/rm/kubernetesrm.
- https://docs.determined.ai/latest/setup-cluster/k8s/internal-task-gateway.html for internal task gateway.
