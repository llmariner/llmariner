#!/usr/bin/env bash

set -euo pipefail

basedir=$(dirname "$0")

# Follow https://docs.konghq.com/kubernetes-ingress-controller/latest/get-started/
#
# The gateway API needs to be installed before Kong intallation as the Kong's helm chart behaves differently based on the presence of the gateway API
# (e.g., whether the cluster role includes HTTPRoutes).
#
# Use the experimental channel to install TCPRoute
kubectl apply -f https://github.com/kubernetes-sigs/gateway-api/releases/download/v1.2.0/experimental-install.yaml

helm repo add kong https://charts.konghq.com
helm repo update
helm upgrade --install kong kong/ingress  -f ./kong_values.yaml

kubectl apply -f "${basedir}"/gateway.yaml
