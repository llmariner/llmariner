#!/usr/bin/env bash

set -euo pipefail

basedir=$(dirname "$0")

helm upgrade \
  --install \
  -n llmariner \
  llmariner \
  oci://public.ecr.aws/cloudnatix/llmariner-charts/llmariner \
  -f "${basedir}"/../llmariner-values.yaml \
  -f "${basedir}"/../llmariner-values-cpu-only.yaml \
  -f "${basedir}"/llmariner-values-llmo-dev.yaml
