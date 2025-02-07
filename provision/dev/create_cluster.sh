#!/usr/bin/env bash

set -eo pipefail

if [ $# -eq 0 ]; then
  echo "Usage: $0 single|multi"
  exit 1
fi

basedir=$(dirname "$0")

case $1 in
  single)
    kind create cluster --name "llmariner-demo" --config "${basedir}"/kind/standalone.yaml
   ;;
  multi)
    kind create cluster --name llmariner-control-plane --config "${basedir}"/kind/control_plane.yaml
    kind create cluster --name llmariner-worker-plane --config "${basedir}"/kind/worker_plane.yaml
    ;;
  gpu-federation)
    kind create cluster --name llmariner-control-plane --config "${basedir}"/kind/control_plane.yaml
    kind create cluster --name llmariner-worker-plane1 --config "${basedir}"/kind/worker_plane.yaml
    kind create cluster --name llmariner-worker-plane2 --config "${basedir}"/kind/worker_plane.yaml
    kind create cluster --name tenant-control-plane --config "${basedir}"/kind/tenant_control_plane.yaml
    ;;
  *)
    echo "Invalid option. Please use 'single' or 'multi'."
    exit 1
    ;;
esac
