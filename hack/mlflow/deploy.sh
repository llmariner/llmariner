#!/usr/bin/env bash

set -euo pipefail
trap 'kill $(jobs -p)' EXIT

dbs=("mlflow" "mlflow_auth")
for db in "${dbs[@]}"; do
  kubectl exec  -n postgres deploy/postgres -- psql -h localhost -U ps_user --no-password -p 5432 -d ps_db -c "CREATE DATABASE ${db};"
done

export AWS_ACCESS_KEY_ID=llmariner-key
export AWS_SECRET_ACCESS_KEY=llmariner-secret

kubectl create secret generic -n mlflow aws \
  --from-literal=accessKeyId=${AWS_ACCESS_KEY_ID} \
  --from-literal=secretAccessKey=${AWS_SECRET_ACCESS_KEY}

kubectl port-forward -n minio service/minio 9000 9090 &
sleep 1

bucket_name=mlflow
aws --endpoint-url http://localhost:9000 s3 mb s3://${bucket_name}

# There is no official Helm chart for MLflow, so we can use the community chart or bitnamicharts/mlflow.
# The community chart hits a bug (https://github.com/community-charts/helm-charts/issues/46), so we use
# bitnamicharts/mlflow here.
#
# See https://github.com/mlflow/mlflow/issues/6118 and https://github.com/bitnami/charts/tree/main/bitnami/mlflow.

helm upgrade \
  --install \
  --create-namespace \
  -n mlflow \
  mlflow oci://registry-1.docker.io/bitnamicharts/mlflow \
  -f values.yaml
