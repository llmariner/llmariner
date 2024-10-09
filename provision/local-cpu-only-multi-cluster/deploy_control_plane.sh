#!/usr/bin/env bash

set -euo pipefail

basedir=$(dirname "$0")

kubectl create namespace llmariner

"${basedir}"/../local-cpu-only/deploy_kong.sh
"${basedir}"/../local-cpu-only/deploy_postgres.sh
"${basedir}"/../local-cpu-only/deploy_minio.sh
"${basedir}"/../local-cpu-only/deploy_milvus.sh
"${basedir}"/deploy_llmariner_control_plane.sh
