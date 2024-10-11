#!/usr/bin/env bash

set -euo pipefail
trap 'kill $(jobs -p)' EXIT

basedir=$(dirname "$0")


kubectl create namespace minio
kubectl apply --namespace minio -f "${basedir}"/../common/minio.yaml
kubectl wait --timeout=60s --for=condition=ready pod -n minio -l app=minio

kubectl port-forward -n minio service/minio 9000 9090 &
sleep 5

minio_user=minioadmin
minio_password=minioadmin

# Obtain the cookie and store in cookies.txt.
curl \
  http://localhost:9090/api/v1/login \
  --cookie-jar cookies.txt \
  --request POST \
  --header 'Content-Type: application/json' \
  --data "{\"accessKey\": \"${minio_user}\", \"secretKey\": \"${minio_password}\"}"

export AWS_ACCESS_KEY_ID=llmariner-key
export AWS_SECRET_ACCESS_KEY=llmariner-secret
# This is needed when the script runs in a GitHub Action runner.
export AWS_EC2_METADATA_DISABLED=true

# Create a new API key.
curl \
  http://localhost:9090/api/v1/service-account-credentials \
  --cookie cookies.txt \
  --request POST \
  --header "Content-Type: application/json" \
  --data "{\"policy\": \"\", \"accessKey\": \"${AWS_ACCESS_KEY_ID}\", \"secretKey\": \"${AWS_SECRET_ACCESS_KEY}\", \"description\": \"\", \"comment\": \"\", \"name\": \"LLMariner\", \"expiry\": null}"

rm cookies.txt

# Create a new bucket.
bucket_name=llmariner
aws --endpoint-url http://localhost:9000 s3 mb s3://${bucket_name}

# Create secrets.
kubectl create secret generic -n llmariner aws \
  --from-literal=accessKeyId=${AWS_ACCESS_KEY_ID} \
  --from-literal=secretAccessKey=${AWS_SECRET_ACCESS_KEY}
