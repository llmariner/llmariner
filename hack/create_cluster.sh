#!/usr/bin/env bash

set -euo pipefail

basedir=$(dirname "$0")

cluster_name="llmariner-demo"

kind create cluster --name "${cluster_name}" --config "${basedir}"/kind-cluster.yaml
