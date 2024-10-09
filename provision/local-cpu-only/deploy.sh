#!/usr/bin/env bash

set -euo pipefail

basedir=$(dirname "$0")

kubectl create namespace llmariner

"${basedir}"/deploy_fake_gpu_operator.sh

"${basedir}"/deploy_kong.sh

"${basedir}"/deploy_postgres.sh

"${basedir}"/deploy_minio.sh

"${basedir}"/deploy_monitoring.sh

"${basedir}"/deploy_milvus.sh

"${basedir}"/deploy_llmariner.sh
