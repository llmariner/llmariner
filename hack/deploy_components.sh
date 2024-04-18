#!/usr/bin/env bash

set -euo pipefail

cluster_name="llm-operator-demo"

model_manager_repo="../../model-manager"
file_manager_repo="../../file-manager"
inference_manager_repo="../../inference-manager"
job_manager_repo="../../job-manager"

# TODO(kenji): This assumes that the HuggingFace API key is stored in the following env var.
kubectl create secret generic -n model-manager hugging-face \
  --from-literal=apiKey="${HUGGING_FACE_HUB_TOKEN}" \

kind load docker-image llm-operator/model-manager-server:latest -n "${cluster_name}"
kind load docker-image llm-operator/model-manager-loader:latest -n "${cluster_name}"
kind load docker-image llm-operator/file-manager-server:latest -n "${cluster_name}"
kind load docker-image llm-operator/inference-manager-engine:latest -n "${cluster_name}"
kind load docker-image llm-operator/job-manager-server:latest -n "${cluster_name}"
kind load docker-image llm-operator/job-manager-dispatcher:latest -n "${cluster_name}"
kind load docker-image llm-operator/experiments-fine-tuning:latest -n "${cluster_name}"

helm upgrade \
  --install \
  -n model-manager \
  model-manager-server \
  "${model_manager_repo}"/deployments/server \
  -f "${model_manager_repo}"/deployments/server/values.yaml \
  -f model-manager-server-values.yaml

helm upgrade \
  --install \
  -n model-manager \
  model-manager-loader \
  "${model_manager_repo}"/deployments/loader \
  -f "${model_manager_repo}"/deployments/loader/values.yaml \
  -f model-manager-loader-values.yaml

helm upgrade \
  --install \
  -n file-manager \
  file-manager-server \
  "${file_manager_repo}"/deployments/server \
  -f "${file_manager_repo}"/deployments/server/values.yaml \
  -f file-manager-server-values.yaml

helm upgrade \
  --install \
  -n inference-manager \
  inference-manager-engine \
  "${inference_manager_repo}"/deployments/engine \
  -f "${inference_manager_repo}"/deployments/engine/values.yaml \
  -f inference-manager-engine-values.yaml

helm upgrade \
  --install \
  -n job-manager \
  job-manager-server \
  "${job_manager_repo}"/deployments/server \
  -f "${job_manager_repo}"/deployments/server/values.yaml \
  -f job-manager-server-values.yaml

helm upgrade \
  --install \
  -n job-manager \
  job-manager-dispatcher \
  "${job_manager_repo}"/deployments/dispatcher \
  -f "${job_manager_repo}"/deployments/dispatcher/values.yaml \
  -f job-manager-dispatcher-values.yaml