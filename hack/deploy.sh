#!/usr/bin/env bash

set -euo pipefail

cluster_name="llm-operator-demo"
job_manager_repo="../../job-manager"

kubectl create namespace inference-server
kind load docker-image inference-server:latest -n "${cluster_name}"
kubectl apply --namespace inference-server -f inference-server.yaml

kubectl create namespace postgres
kubectl apply --namespace postgres -f postgres.yaml
kubectl exec  -n postgres deploy/postgres -- psql -h localhost -U ps_user --no-password -p 5432 -d ps_db -c "CREATE DATABASE job_manager;"

kubectl create namespace job-manager
kubectl apply -n job-manager -f job-manager-postgres-secret.yaml
kind load docker-image job-manager-server:latest -n "${cluster_name}"
helm upgrade \
  --install \
  -n job-manager \
  job-manager-server \
  "${job_manager_repo}"/deployments/server \
  -f "${job_manager_repo}"/deployments/server/values.yaml \
  -f job-manager-server-values.yaml
kind load docker-image job-manager-dispatcher:latest -n "${cluster_name}"
helm upgrade \
  --install \
  -n job-manager \
  job-manager-dispatcher \
  "${job_manager_repo}"/deployments/dispatcher \
  -f "${job_manager_repo}"/deployments/dispatcher/values.yaml \
  -f job-manager-dispatcher-values.yaml
