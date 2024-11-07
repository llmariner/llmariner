#!/usr/bin/env bash

set -eo pipefail

ACTION=${1:?ACTION}
KUBECONFIG_CONTEXT=${2:?KUBECONFIG_CONTEXT}
NAMESPACE=${3:?NAMESPACE}

basedir=$(dirname "$0")

os=$(uname -s)

if [ "${os}" == "Linux" ]; then
  kubectl ${ACTION} \
          --context="${KUBECONFIG_CONTEXT}" \
          --namespace="${NAMESPACE}" \
          -f "${basedir}/control_plane_service_linux.yaml"
else
  kubectl ${ACTION} \
          --context="${KUBECONFIG_CONTEXT}" \
          --namespace="${NAMESPACE}" \
          -f "${basedir}/control_plane_service.yaml"
fi
