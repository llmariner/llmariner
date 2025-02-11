# Provision LLMariner for development

This directory contains scripts that create a Kind cluster locally and deploy LLMariner and other necessary components.

We use [Fake GPU operator](https://github.com/run-ai/fake-gpu-operator) so that LLMariner can run in a machine that doesn't have any GPU. Ollama (not vLLM) is used as inference runtime.

## Requirements

- [Docker](https://docs.docker.com/engine/install/)
- [Kind](https://kind.sigs.k8s.io/docs/user/quick-start/#installation)
- [Helmfile](https://helmfile.readthedocs.io/en/latest/#installation)

## Provisioning

### Standalone Mode

```bash
./create_cluster.sh single
helmfile init
helmfile apply --skip-diff-on-install
```

> [!TIP]
> You can filter the components to deploy using the `--selector(-l)` flag.
> For example, to filter out the monitoring components, set the `-l tier!=monitoring` flag.
> For deploying just the llmariner, use `-l app=llmariner`.

### Multi-Cluster Mode

#### Single Worker Cluster

```bash
./create_cluster.sh multi
helmfile init
helmfile apply -e control -l app!=fake-gpu-operator,tier!=monitoring --skip-diff-on-install

# Please set the endpoint address to http://localhost/v1
llma auth login
export REGISTRATION_KEY=$(llma admin clusters register worker-cluster | sed -n 's/.*Registration Key: "\([^"]*\)".*/\1/p')
helmfile apply -e worker -l app=fake-gpu-operator -l tier=monitoring -l app=llmariner --skip-diff-on-install
```

> [!NOTE]
> The worker cluster uses an ExternalName service to reach the control plane.
> Please note that the current service definition is for Mac/Windows (Docker Desktop).
> See https://github.com/kubernetes-sigs/kind/issues/1200#issuecomment-130485579.

> [!NOTE]
> Please note that the endpoint address is http://localhost/v1, not http://localhost:8080/v1.

#### For GPU Federation: Multiple Worker Cluster with Tenant Control-Plane cluster

```bash
./create_cluster.sh gpu-federation
helmfile init
helmfile apply -e control -l app!=fake-gpu-operator,tier!=monitoring --skip-diff-on-install

# Please set the endpoint address to http://localhost/v1
llma auth login

# Unregister the default cluster.
llma admin clusters unregister Default

# Deploy two worker clusters.
export REGISTRATION_KEY=$(llma admin clusters register worker-cluster1 | sed -n 's/.*Registration Key: "\([^"]*\)".*/\1/p')
helmfile apply -e worker --kube-context kind-llmariner-worker-plane1 -l app=fake-gpu-operator -l app=llmariner --skip-diff-on-install

export REGISTRATION_KEY=$(llma admin clusters register worker-cluster2 | sed -n 's/.*Registration Key: "\([^"]*\)".*/\1/p')
helmfile apply -e worker --kube-context kind-llmariner-worker-plane2 -l app=fake-gpu-operator -l app=llmariner --skip-diff-on-install

# Deploy to a tenant control-plane cluster.
export TENANT_API_KEY=$(llma auth api-keys create tenant -o 'Default Organization' --role tenant-system --service-account | sed -n 's/.*Secret: \(.*\)/\1/p')
helmfile apply -e tenant-control -l app=llmariner --skip-diff-on-install
```

> [!NOTE]
> The worker cluster uses an ExternalName service to reach the control plane.
> Please note that the current service definition is for Mac/Windows (Docker Desktop).
> See https://github.com/kubernetes-sigs/kind/issues/1200#issuecomment-130485579.

> [!NOTE]
> Please note that the endpoint address is http://localhost/v1, not http://localhost:8080/v1.


## Testing

```bash
LLMARINER_API_KEY=default-key-secret ./validate_deployment.sh
```
