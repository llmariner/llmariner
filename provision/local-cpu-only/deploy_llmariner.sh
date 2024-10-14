#!/usr/bin/env bash

set -euo pipefail

: ${CHART_LOCATION:=oci://public.ecr.aws/cloudnatix/llmariner-charts/llmariner}

basedir=$(dirname "$0")

helm upgrade \
  --install \
  -n llmariner \
  llmariner \
  "${CHART_LOCATION}" \
  -f "${basedir}"/../common/llmariner-values.yaml
