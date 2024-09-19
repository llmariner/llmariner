#!/usr/bin/env bash

set -euo pipefail

basedir=$(dirname "$0")

kubectl create namespace postgres
kubectl apply --namespace postgres -f "${basedir}"/postgres.yaml
kubectl apply -n llm-operator -f "${basedir}"/postgres-secret.yaml
