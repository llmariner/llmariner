#!/usr/bin/env bash

set -euo pipefail
trap 'kill $(jobs -p)' EXIT

kubectl create namespace minio
kubectl apply --namespace minio -f minio.yaml
kubectl wait --for=condition=ready pod -n minio -l app=minio

kubectl port-forward -n minio service/minio 9000 9090 &
sleep 5

minio_user=minioadmin
minio_password=minioadmin

curl \
  http://localhost:9090/api/v1/login \
  --cookie-jar cookies.txt \
  --request POST \
  --header 'Content-Type: application/json' \
  --data "{\"accessKey\": \"${minio_user}\", \"secretKey\": \"${minio_password}\"}"

export AWS_ACCESS_KEY_ID=llm-operator-key
export AWS_SECRET_ACCESS_KEY=llm-operator-secret

curl \
  http://localhost:9090/api/v1/service-account-credentials \
  --cookie cookies.txt \
  --request POST \
  --header "Content-Type: application/json" \
  --data "{\"policy\": \"\", \"accessKey\": \"${AWS_ACCESS_KEY_ID}\", \"secretKey\": \"${AWS_SECRET_ACCESS_KEY}\", \"description\": \"\", \"comment\": \"\", \"name\": \"LLM Operator\", \"expiry\": null}"

rm cookies.txt

bucket_name=llm-operator
aws --endpoint-url http://localhost:9000 s3 mb s3://${bucket_name}

kubectl create secret generic -n file-manager aws \
  --from-literal=accessKeyId=${AWS_ACCESS_KEY_ID} \
  --from-literal=secretAccessKey=${AWS_SECRET_ACCESS_KEY}

kubectl create secret generic -n job-manager aws \
  --from-literal=accessKeyId=${AWS_ACCESS_KEY_ID} \
  --from-literal=secretAccessKey=${AWS_SECRET_ACCESS_KEY}

kubectl create secret generic -n inference-manager aws \
  --from-literal=accessKeyId=${AWS_ACCESS_KEY_ID} \
  --from-literal=secretAccessKey=${AWS_SECRET_ACCESS_KEY}
