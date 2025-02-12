#! /usr/bin/env bash

set -xe

llma auth api-keys delete tenant || true

for cluster in "gpu-worker-cluster-large gpu-worker-cluster-small"; do
  llma admin clusters unregister "${cluster}" || true
done

kind delete clusters tenant-cluster gpu-worker-cluster-large gpu-worker-cluster-small || true
