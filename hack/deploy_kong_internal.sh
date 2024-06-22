#!/usr/bin/env bash

set -euo pipefail

basedir=$(dirname "$0")

helm repo add kong https://charts.konghq.com
helm repo update
helm install --create-namespace kong-proxy kong/kong -n kong-internal --set ingressController.installCRDs=false -f "${basedir}"/kong_internal_values.yaml
