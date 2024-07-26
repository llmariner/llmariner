#!/usr/bin/env bash

set -euo pipefail

basedir=$(dirname "$0")

# Follow https://docs.konghq.com/kubernetes-ingress-controller/latest/get-started/
#
# The gateway API needs to be installed before Kong intallation as the Kong's helm chart behaves differently based on the presence of the gateway API
# (e.g., whether the cluster role includes HTTPRoutes).
kubectl apply -f https://github.com/kubernetes-sigs/gateway-api/releases/download/v1.1.0/standard-install.yaml

helm repo add kong https://charts.konghq.com
helm repo update
helm upgrade --install --create-namespace kong kong/ingress -n llm-operator -f ./kong_values.yaml

kubectl apply -f "${basedir}"/gateway.yaml
