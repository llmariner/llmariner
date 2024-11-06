#!/usr/bin/env bash

set -eo pipefail

KUBECONFIG_CONTEXT=${1}

basedir=$(dirname "$0")

os=$(uname -s)

if [ "${os}" == "Linux" ]; then
  kubectl --context="${KUBECONFIG_CONTEXT}" apply -f "${basedir}/control_plane_service_linux.yaml"
else
  kubectl --context="${KUBECONFIG_CONTEXT}" apply -f "${basedir}/control_plane_service.yaml"
fi
