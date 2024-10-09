#!/usr/bin/env bash

set -euo pipefail

basedir=$(dirname "$0")

kind create cluster --name llmariner-control-plane --config "${basedir}"/kind_cluster_control_plane.yaml

kind create cluster --name llmariner-worker-plane --config "${basedir}"/kind_cluster_worker_plane.yaml
