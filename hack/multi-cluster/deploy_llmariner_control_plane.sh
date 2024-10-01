#!/usr/bin/env bash

set -euo pipefail

basedir=$(dirname "$0")

helm upgrade \
  --install \
  --wait \
  -n llmariner \
  llmariner \
  oci://public.ecr.aws/cloudnatix/llmariner-charts/llmariner \
  --set tags.worker=false \
  -f "${basedir}"/../llmariner-values.yaml \
  -f "${basedir}"/../llmariner-values-cpu-only.yaml \
  -f "${basedir}"/llmariner-values-control-plane.yaml
