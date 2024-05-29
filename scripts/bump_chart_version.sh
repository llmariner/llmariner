#!/usr/bin/env bash
#
# Print the latest version of each LLM component.
#
# TODO(kenji): Update Chart.yaml with the latest version.

set -euo pipefail

repos=(
    "file-manager"
    "inference-manager"
    "job-manager"
    "model-manager"
    "user-manager"
    "rbac-manager"
    "session-manager"
)

for repo in "${repos[@]}"; do
   latest_version=$(git ls-remote --tags --refs --sort="v:refname" "https://github.com/llm-operator/${repo}.git" | tail -n 1 | sed 's/.*\///')
   echo $repo $latest_version
done
