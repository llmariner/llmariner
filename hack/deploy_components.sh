#!/usr/bin/env bash

set -euo pipefail

cluster_name="llm-operator-demo"

llm_operator_namespace=llm-operator

helm repo add llm-operator http://llm-operator-charts.s3-website-us-west-2.amazonaws.com/
helm repo update

helm upgrade \
  --install \
  -n "${llm_operator_namespace}" \
  model-manager-server \
  llm-operator/model-manager-server \
  -f model-manager-server-values.yaml

helm upgrade \
  --install \
  -n "${llm_operator_namespace}" \
  model-manager-loader \
  llm-operator/model-manager-loader \
  -f model-manager-loader-values.yaml

helm upgrade \
  --install \
  -n "${llm_operator_namespace}" \
  file-manager-server \
  llm-operator/file-manager-server \
  -f file-manager-server-values.yaml

helm upgrade \
  --install \
  -n "${llm_operator_namespace}" \
  inference-manager-engine \
  llm-operator/inference-manager-engine \
  -f inference-manager-engine-values.yaml

helm upgrade \
  --install \
  -n "${llm_operator_namespace}" \
  job-manager-server \
  llm-operator/job-manager-server \
  -f job-manager-server-values.yaml

helm upgrade \
  --install \
  -n "${llm_operator_namespace}" \
  job-manager-dispatcher \
  llm-operator/job-manager-dispatcher \
  -f job-manager-dispatcher-values.yaml
