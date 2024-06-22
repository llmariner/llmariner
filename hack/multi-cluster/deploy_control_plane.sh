#!/usr/bin/env bash

set -euo pipefail

basedir=$(dirname "$0")

kubectl create namespace llm-operator

"${basedir}"/../deploy_kong.sh
"${basedir}"/../deploy_postgres.sh
"${basedir}"/../deploy_minio.sh
"${basedir}"/../deploy_milvus.sh
"${basedir}"/../deploy_keda.sh
"${basedir}"/deploy_llm_operator_control_plane.sh

kubectl apply -n llm-operator -f "${basedir}"/session_manager_server_service.yaml

# Create a service of node port to map port 80 to session-manager-server
