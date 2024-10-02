#!/usr/bin/env bash

set -euo pipefail

basedir=$(dirname "$0")

helm upgrade \
  --install \
  --wait \
  -n llmariner-wp \
  llmariner \
  oci://public.ecr.aws/cloudnatix/llmariner-charts/llmariner \
  -f "${basedir}"/../llmariner-values.yaml \
  -f "${basedir}"/llmariner-values-worker-plane.yaml
