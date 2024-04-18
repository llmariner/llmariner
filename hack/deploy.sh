#!/usr/bin/env bash

set -euo pipefail

cluster_name="llm-operator-demo"

# TODO(kenji): Consider running all components in the same namespace to simplify the deployment.
kubectl create namespace model-manager
kubectl create namespace file-manager
kubectl create namespace inference-manager
kubectl create namespace job-manager

./deploy_fake_gpu_operator.sh

./deploy_kong.sh

./deploy_postgres.sh

./deploy_minio.sh

./deploy_components.sh
