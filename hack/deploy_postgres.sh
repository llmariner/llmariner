#!/usr/bin/env bash

set -euo pipefail

kubectl create namespace postgres
kubectl apply --namespace postgres -f postgres.yaml
kubectl wait --timeout=60s --for=condition=ready pod -n postgres -l app=postgres
# Wait for extra seconds
sleep 5

dbs=("model_manager" "file_manager" "job_manager", "dex")
for db in "${dbs[@]}"; do
  kubectl exec  -n postgres deploy/postgres -- psql -h localhost -U ps_user --no-password -p 5432 -d ps_db -c "CREATE DATABASE ${db};"
done

kubectl apply -n llm-operator -f postgres-secret.yaml
