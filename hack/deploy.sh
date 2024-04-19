#!/usr/bin/env bash

set -euo pipefail

cluster_name="llm-operator-demo"

kubectl create namespace llm-operator
kubectl create namespace llm-operator-jobs

./deploy_fake_gpu_operator.sh

./deploy_kong.sh

./deploy_postgres.sh

./deploy_minio.sh

./deploy_components.sh
