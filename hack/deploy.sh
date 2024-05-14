#!/usr/bin/env bash

set -euo pipefail

kubectl create namespace llm-operator
# This namespace is for fine-tuning jobs
kubectl create namespace example-org

./deploy_fake_gpu_operator.sh

./deploy_kong.sh

./deploy_postgres.sh

./deploy_minio.sh

./deploy_monitoring.sh

./deploy_keda.sh

./deploy_components.sh
