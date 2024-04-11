#!/usr/bin/env bash

set -euo pipefail

./deploy_kong.sh

cluster_name="llm-operator-demo"

model_manager_repo="../../model-manager"
file_manager_repo="../../file-manager"
inference_manager_repo="../../inference-manager"
job_manager_repo="../../job-manager"


kubectl create namespace postgres
kubectl create namespace minio
kubectl create namespace model-manager
kubectl create namespace file-manager
kubectl create namespace inference-manager
kubectl create namespace job-manager

kubectl apply --namespace postgres -f postgres.yaml

kubectl apply --namespace minio -f minio.yaml

# TODO(kenji): Run this after the postgres pod starts running.
kubectl exec  -n postgres deploy/postgres -- psql -h localhost -U ps_user --no-password -p 5432 -d ps_db -c "CREATE DATABASE model_manager;"
kubectl exec  -n postgres deploy/postgres -- psql -h localhost -U ps_user --no-password -p 5432 -d ps_db -c "CREATE DATABASE file_manager;"
kubectl exec  -n postgres deploy/postgres -- psql -h localhost -U ps_user --no-password -p 5432 -d ps_db -c "CREATE DATABASE job_manager;"

kubectl apply -n model-manager -f postgres-secret.yaml
kubectl apply -n file-manager -f postgres-secret.yaml
kubectl apply -n job-manager -f postgres-secret.yaml

kubectl apply -f model-store.yaml

kind load docker-image llm-operator/model-manager-server:latest -n "${cluster_name}"
kind load docker-image llm-operator/file-manager-server:latest -n "${cluster_name}"
kind load docker-image llm-operator/inference-manager-engine:latest -n "${cluster_name}"
kind load docker-image llm-operator/job-manager-server:latest -n "${cluster_name}"
kind load docker-image llm-operator/job-manager-dispatcher:latest -n "${cluster_name}"
# kind load docker-image llm-operator/experiments-fine-tuning:latest -n "${cluster_name}"
kind load docker-image llm-operator/experiments-fake-job:latest -n "${cluster_name}"

helm upgrade \
  --install \
  -n model-manager \
  model-manager-server \
  "${model_manager_repo}"/deployments/server \
  -f "${model_manager_repo}"/deployments/server/values.yaml \
  -f model-manager-server-values.yaml

helm upgrade \
  --install \
  -n file-manager \
  file-manager-server \
  "${file_manager_repo}"/deployments/server \
  -f "${file_manager_repo}"/deployments/server/values.yaml \
  -f file-manager-server-values.yaml

helm upgrade \
  --install \
  -n inference-manager \
  inference-manager-engine \
  "${inference_manager_repo}"/deployments/engine \
  -f "${inference_manager_repo}"/deployments/engine/values.yaml \
  -f inference-manager-engine-values.yaml

helm upgrade \
  --install \
  -n job-manager \
  job-manager-server \
  "${job_manager_repo}"/deployments/server \
  -f "${job_manager_repo}"/deployments/server/values.yaml \
  -f job-manager-server-values.yaml

helm upgrade \
  --install \
  -n job-manager \
  job-manager-dispatcher \
  "${job_manager_repo}"/deployments/dispatcher \
  -f "${job_manager_repo}"/deployments/dispatcher/values.yaml \
  -f job-manager-dispatcher-values.yaml
