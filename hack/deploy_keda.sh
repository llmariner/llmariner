#!/usr/bin/env bash

set -euo pipefail

helm repo add kedacore https://kedacore.github.io/charts
helm repo update
helm upgrade --install --wait \
    --namespace keda  \
    --create-namespace \
    keda kedacore/keda
