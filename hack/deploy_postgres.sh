#!/usr/bin/env bash

set -euo pipefail

kubectl create namespace postgres
kubectl apply --namespace postgres -f postgres.yaml
kubectl wait --timeout=60s --for=condition=ready pod -n postgres -l app=postgres
# Wait for extra seconds
sleep 5
kubectl exec  -n postgres deploy/postgres -- psql -h localhost -U ps_user --no-password -p 5432 -d ps_db -c "CREATE DATABASE model_manager;"
kubectl exec  -n postgres deploy/postgres -- psql -h localhost -U ps_user --no-password -p 5432 -d ps_db -c "CREATE DATABASE file_manager;"
kubectl exec  -n postgres deploy/postgres -- psql -h localhost -U ps_user --no-password -p 5432 -d ps_db -c "CREATE DATABASE job_manager;"

kubectl apply -n model-manager -f postgres-secret.yaml
kubectl apply -n file-manager -f postgres-secret.yaml
kubectl apply -n job-manager -f postgres-secret.yaml
