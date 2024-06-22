#!/usr/bin/env bash

set -euo pipefail

basedir=$(dirname "$0")

kind create cluster --name llm-operator-control-plane --config "${basedir}"/kind_cluster_control_plane.yaml

kind create cluster --name llm-operator-worker-plane --config "${basedir}"/kind_cluster_worker_plane.yaml
