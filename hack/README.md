# Test scripts

```bash
./create_cluster.sh
./deploy.sh
./run-test.sh
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
  database: job_manager
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
    cacheDir: /Users/kenji/.cache/hugging-face-cache/hub
```

Then you can set up port-forwarding and run `loader`.

```bash
kubectl port-forward -n postgres service/postgres 5432 &
kubectl port-forward -n minio service/minio 9000 9090 &

export AWS_ACCESS_KEY_ID=llm-operator-key
export AWS_SECRET_ACCESS_KEY=llm-operator-secret

$(job_manager_repo)/bin/loader run --config config.yaml
```
