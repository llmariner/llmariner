#!/usr/bin/env bash

set -euo pipefail

helm repo add milvus https://zilliztech.github.io/milvus-helm/
helm repo update
helm upgrade --install milvus milvus/milvus --namespace milvus --create-namespace -f milvus_values.yaml
