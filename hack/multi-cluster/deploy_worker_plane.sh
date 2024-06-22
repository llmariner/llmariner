#!/usr/bin/env bash

set -euo pipefail

basedir=$(dirname "$0")

kubectl create namespace llm-operator-wp

"${basedir}"/../deploy_fake_gpu_operator.sh
"${basedir}"/../deploy_kong_internal.sh

export AWS_ACCESS_KEY_ID=llm-operator-key
export AWS_SECRET_ACCESS_KEY=llm-operator-secret
kubectl create secret generic -n llm-operator-wp aws \
  --from-literal=accessKeyId=${AWS_ACCESS_KEY_ID} \
  --from-literal=secretAccessKey=${AWS_SECRET_ACCESS_KEY}

kubectl apply -n llm-operator-wp -f "${basedir}"/control_plane_service.yaml

# Create a cluster registration credential
REGISTRATION_KEY=$(llmo admin clusters register worker-cluster | sed -n 's/.*Registration Key: "\([^"]*\)".*/\1/p')
kubectl create secret generic \
  -n llm-operator-wp \
  cluster-registration-key \
  --from-literal=regKey="${REGISTRATION_KEY}"

"${basedir}"/deploy_llm_operator_worker_plane.sh
