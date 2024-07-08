#!/usr/bin/env bash

set -euo pipefail

basedir=$(dirname "$0")

helm upgrade \
  --install \
  -n llm-operator \
  llm-operator \
  oci://public.ecr.aws/cloudnatix/llm-operator-charts/llm-operator \
  -f "${basedir}"/../llm-operator-values.yaml \
  -f "${basedir}"/../llm-operator-values-cpu-only.yaml \
  -f "${basedir}"/llm-operator-values-llmo-dev.yaml
