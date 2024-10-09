# Test scripts

## Build a Kind cluster and deploy LLMariner in a single non-GPU node

If you want to use a Helm chart in your local filesystem, update `deployments/llmariner/Chart.yaml`
and specify the Helm chart location:

```yaml
- name: model-manager-server
  repository: "file://../../../model-manager/deployments/server"
  version: "*"
```

If you also want to use a different container image, add the following to `llmariner-values.yaml`.

```yaml
model-manager-server:
  image:
    repository: llmariner/model-manager-server
    pullPolicy: Never
  version: latest
```

Then load the image to the Kind cluster and deploy LLMariner.

```bash
kind load docker-image llmariner/model-manager-server:latest -n llmariner-demo

helm dependencies build deployments/llmariner
helm upgrade --install -n llmariner llmariner ./deployments/llmariner  -f hack/llmariner-values.yaml
```

## Run a Fake Fine-Tuning Job

Running a fine-tuning job requires GPU. If you want to test an end-to-end flow of fine-tuning without GPU, you can configure
`job-manager-dispatcher` to create a fake job by adding the following to `values.yaml`:

```yaml
job-manager-dispatcher:
  job:
    image: public.ecr.aws/cloudnatix/llmariner/fake-job
    version: latest
    imagePullPolicy: IfNotPresent
```

The fake-job already has an output model in its container image, and it just copies to the output directory.

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
    bucket: llmariner
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

export AWS_ACCESS_KEY_ID=llmariner-key
export AWS_SECRET_ACCESS_KEY=llmariner-secret

$(job_manager_repo)/bin/loader run --config config.yaml
```
