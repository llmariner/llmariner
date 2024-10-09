#!/usr/bin/env bash

set -euo pipefail

basedir=$(dirname "$0")

kubectl create namespace postgres
kubectl apply --namespace postgres -f "${basedir}"/../common/postgres.yaml
kubectl apply -n llmariner -f "${basedir}"/postgres-secret.yaml
