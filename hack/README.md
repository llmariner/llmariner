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
