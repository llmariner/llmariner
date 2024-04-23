#!/usr/bin/env bash

set -euo pipefail

kubectl create namespace llm-operator
kubectl create namespace llm-operator-jobs

./deploy_fake_gpu_operator.sh

./deploy_kong.sh

./deploy_postgres.sh

./deploy_minio.sh

./deploy_components.sh
