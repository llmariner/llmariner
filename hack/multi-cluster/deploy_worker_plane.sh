#!/usr/bin/env bash

set -euo pipefail

basedir=$(dirname "$0")

kubectl create namespace llmariner-wp

"${basedir}"/../deploy_fake_gpu_operator.sh

export AWS_ACCESS_KEY_ID=llmariner-key
export AWS_SECRET_ACCESS_KEY=llmariner-secret
kubectl create secret generic -n llmariner-wp aws \
  --from-literal=accessKeyId=${AWS_ACCESS_KEY_ID} \
  --from-literal=secretAccessKey=${AWS_SECRET_ACCESS_KEY}

kubectl apply -n llmariner-wp -f "${basedir}"/control_plane_service.yaml

# Create a cluster registration credential
REGISTRATION_KEY=$(llma admin clusters register worker-cluster | sed -n 's/.*Registration Key: "\([^"]*\)".*/\1/p')
kubectl create secret generic \
  -n llmariner-wp \
  cluster-registration-key \
  --from-literal=regKey="${REGISTRATION_KEY}"

"${basedir}"/deploy_llmariner_worker_plane.sh
