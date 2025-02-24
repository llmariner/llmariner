#!/usr/bin/env bash
#
# Fetch or update the repositories.

set -euo pipefail

basedir=$(dirname "$0")

mkdir -p "${basedir}"/repos
cd "${basedir}"/repos

repos="api-usage cluster-manager file-manager inference-manager job-manager model-manager rbac-manager session-manager user-manager vector-store-manager"

for repo in $repos; do
  if [ ! -d "${repo}" ]; then
      git clone https://github.com/llmariner/${repo}.git
  else
      cd "${repo}"
      git checkout main
      git pull
  fi
done
