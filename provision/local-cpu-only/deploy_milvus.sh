#!/usr/bin/env bash

set -euo pipefail

basedir=$(dirname "$0")

helm repo add milvus https://zilliztech.github.io/milvus-helm/
helm repo update
helm upgrade --install milvus milvus/milvus --namespace milvus --create-namespace -f "${basedir}"/../common/milvus_values.yaml
