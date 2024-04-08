#!/usr/bin/env bash

set -euo pipefail

./deploy_kong.sh

cluster_name="llm-operator-demo"
inference_manager_repo="../../inference-manager"
job_manager_repo="../../job-manager"

kubectl create namespace postgres
kubectl create namespace inference-manager
kubectl create namespace job-manager
kubectl create namespace model-store

kubectl apply --namespace postgres -f postgres.yaml
# TODO(kenji): Run this after the postgres pod starts running.
kubectl exec  -n postgres deploy/postgres -- psql -h localhost -U ps_user --no-password -p 5432 -d ps_db -c "CREATE DATABASE job_manager;"

kubectl apply -n job-manager -f job-manager-postgres-secret.yaml

kubectl apply -f model-store.yaml

kind load docker-image llm-operator/inference-manager-engine:latest -n "${cluster_name}"
kind load docker-image llm-operator/job-manager-server:latest -n "${cluster_name}"
kind load docker-image llm-operator/job-manager-dispatcher:latest -n "${cluster_name}"
# kind load docker-image llm-operator/experiments-fine-tuning:latest -n "${cluster_name}"

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
