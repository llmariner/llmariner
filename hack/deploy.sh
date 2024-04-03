#!/usr/bin/env bash

set -euo pipefail

cluster_name="llm-operator-demo"

kind load docker-image inference-server:latest -n "${cluster_name}"
kind load docker-image job-manager-server:latest -n "${cluster_name}"

kubectl create namespace inference-server
kubectl apply --namespace inference-server -f inference-server.yaml

kubectl create namespace postgres
kubectl apply --namespace postgres -f postgres.yaml

job_manager_repo="../../job-manager"

kubectl create namespace job-manager
helm upgrade --install -n job-manager job-manager-server "${job_manager_repo}"/deployments/server
helm upgrade --install -n job-manager job-manager-dispatcher "${job_manager_repo}"/deployments/dispatcher
