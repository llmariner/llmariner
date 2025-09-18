#!/usr/bin/env bash

set -eo pipefail
trap 'kill $(jobs -p)' EXIT

MINIO_USER=${1:?MINIO_USER}
MINIO_PASS=${2:?MINIO_PASS}
ACCESS_KEY=${3:?ACCESS_KEY}
SECRET_KEY=${4:?SECRET_KEY}
KUBECONFIG_CONTEXT=${5}
NAMESPACE=${6:?NAMESPACE}

for i in {1..100}; do
  kubectl get pod --context="${KUBECONFIG_CONTEXT}" --namespace="${NAMESPACE}" -l app.kubernetes.io/name=minio
  sleep 1
done

kubectl get pod --context="${KUBECONFIG_CONTEXT}" --namespace="${NAMESPACE}" -l app.kubernetes.io/name=minio -o yaml

kubectl wait pod \
        --context="${KUBECONFIG_CONTEXT}" \
        --timeout=180s \
        --for=condition=ready \
        --namespace="${NAMESPACE}" \
        -l app.kubernetes.io/name=minio

kubectl port-forward \
        --context="${KUBECONFIG_CONTEXT}" \
        --namespace="${NAMESPACE}" \
        service/minio 9001 &

sleep 5

# Obtain the cookie and store in cookies.txt.
curl -fs --show-error http://localhost:9001/api/v1/login \
  --cookie-jar cookies.txt \
  --header 'Content-Type: application/json' \
  --data @- << EOF
{
  "accessKey": "$MINIO_USER",
  "secretKey": "$MINIO_PASS"
}
EOF

# Create a new API key.
curl http://localhost:9001/api/v1/service-account-credentials \
  --cookie cookies.txt \
  --header "Content-Type: application/json" \
  --data @- << EOF >/dev/null
{
  "name": "LLMariner",
  "accessKey": "$ACCESS_KEY",
  "secretKey": "$SECRET_KEY",
  "description": "",
  "comment": "",
  "policy": "",
  "expiry": null
}
EOF

rm cookies.txt
