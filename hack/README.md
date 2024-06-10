# Test scripts

## Build a Kind cluster and deploy LLM Operator in a single non-GPU node

Run the following commands:

```bash
./create_cluster.sh
./deploy.sh
./check_readiness.sh
```

If you want to use a Helm chart in your local filesystem, update `deployments/llm-operator/Chart.yaml`
and specify the Helm chart location:

```yaml
- name: model-manager-server
  repository: "file://../../../model-manager/deployments/server"
  version: "*"
```

If you also want to use a different container image, add the following to `llm-operator-values.yaml`.

```yaml
model-manager-server:
  image:
    repository: llm-operator/model-manager-server
    pullPolicy: Never
  version: latest
```

Then load the image to the Kind cluster and deploy LLM Operator.

```bash
kind load docker-image llm-operator/model-manager-server:latest -n llm-operator-demo

helm dependencies build deployments/llm-operator
helm upgrade --install -n llm-operator llm-operator ./deployments/llm-operator  -f hack/llm-operator-values.yaml
```

## Run a Fake Fine-Tuning Job

Running a fine-tuning job requires GPU. If you want to test an end-to-end flow of fine-tuning without GPU, you can configure
`job-manager-dispatcher` to create a fake job by adding the following to `values.yaml`:

```yaml
job-manager-dispatcher:
  job:
    image: public.ecr.aws/v8n3t7y5/llm-operator/fake-job
    version: latest
    imagePullPolicy: IfNotPresent
```

The fake-job already has an output model in its container image, and it just copies to the output directory.

You can also update `deploy_llm_operator.sh` to take `llm-operator-values-cpu-only.yaml` when deploying a Helm chart.

```console
helm upgrade \
  --install \
  -n llm-operator \
  llm-operator \
  llm-operator/llm-operator \
  -f llm-operator-values.yaml \
  -f llm-operator-values-cpu-only.yaml
```

## Deployment to a Nvidia H100 Launchpad Instance

```bash
helm upgrade \
  --install \
  -n llm-operator \
  llm-operator \
  llm-operator/llm-operator \
  -f llm-operator-values.yaml \
  -f llm-operator-values-nvidia-launchpad.yaml
```

## Setting up MinIO

See https://min.io/docs/minio/kubernetes/upstream/index.html

First, set up port-forwarding.

```bash
kubectl port-forward -n minio service/minio 9000 9090
```

Access http://localhost:9090. The username and the password are both `minioadmin`.

Generate an API key.

Then you can make an API call with the AWSCLI:

```bash
export AWS_ACCESS_KEY_ID=...
export AWS_SECRET_ACCESS_KEY=...
aws --endpoint-url http://localhost:9000 s3 mb s3://test-bucket
aws --endpoint-url http://localhost:9000 s3 ls
```

`./deploy_minio.sh` creates an API key so that components can access MinIO. To use the API key,
set the env vars in the following way:

```bash
export AWS_ACCESS_KEY_ID=llm-operator-key
export AWS_SECRET_ACCESS_KEY=llm-operator-secret
aws --endpoint-url http://localhost:9000 s3 ls --recursive s3://llm-operator
```

## Loading base images from your local HuggingFace cache

If you want to load models from your local HuggingFace cache, you can run `model-manager-loader` locally with
following the config:

```yaml
database:
  host: localhost
  port: 5432
  database: model_manager
  username: ps_user
  passwordEnvName: DB_PASSWORD

objectStore:
  s3:
    endpointUrl: http://localhost:9000
    bucket: llm-operator
    pathPrefix: models
    baseModelPathPrefix: base-models

baseModels:
- google/gemma-2b

modelLoadInterval: 1m

downloader:
  huggingFace:
    # Change this to your cache directory.
    cacheDir: /Users/kenji/.cache/huggingface/hub
```

Then you can set up port-forwarding and run `loader`.

```bash
kubectl port-forward -n postgres service/postgres 5432 &
kubectl port-forward -n minio service/minio 9000 9090 &

export AWS_ACCESS_KEY_ID=llm-operator-key
export AWS_SECRET_ACCESS_KEY=llm-operator-secret

$(job_manager_repo)/bin/loader run --config config.yaml
```
