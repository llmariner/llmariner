#!/usr/bin/env bash

set -euo pipefail

if [[ `kubectl exec -n postgres deploy/postgres -- psql -h localhost -U ps_user --no-password -p 5432 -d ps_db -c '\list' | grep 'file_manager' | wc -l` -gt 0 ]]; then
  echo "Databases already exist"
  exit 0
fi

dbs=("cluster_manager" "model_manager" "file_manager" "job_manager" "user_manager" "dex" "vector_store_manager")
for db in "${dbs[@]}"; do
  kubectl exec  -n postgres deploy/postgres -- psql -h localhost -U ps_user --no-password -p 5432 -d ps_db -c "CREATE DATABASE ${db};"
done
