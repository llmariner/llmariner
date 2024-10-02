# Enable Time Slicing in a Kind Cluster with Nvidia GPU Operator

## Overview

This document describes the steps to install Nvidia GPU operator in a KIND cluster and enable time-slicing GPU sharing in the cluster.

## Procedure

### Step 1. Setup KIND cluster with Nvidia GPU

Follow [the instruction](https://github.com/llmariner/llmariner/blob/main/docs/developments/build_kind_cluster_with_gpu.md) to set up a KIND cluster with Nvidia GPU.

### Step 2. Install Nvidia GPU Operator

```console
helm repo add nvidia https://helm.ngc.nvidia.com/nvidia
helm repo update
helm install --wait --generate-name \
    -n gpu-operator --create-namespace \
    nvidia/gpu-operator \
    --set cdi.enabled=true \
    --set driver.enabled=false \
    --set toolkit.enabled=false
```

### Step 3. Configure time-slicing GPU sharing

In this example, configure the cluster with time-slicing based GPU `replicas` to be 4.

```console
$ cat ./time-slicing-config-all.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: time-slicing-config-all
data:
  any: |-
    version: v1
    flags:
      migStrategy: none
    sharing:
      timeSlicing:
        resources:
        - name: nvidia.com/gpu
          replicas: 4

$ kubectl create -n gpu-operator -f time-slicing-config-all.yaml

$ kubectl patch clusterpolicies.nvidia.com/cluster-policy \
    -n gpu-operator --type merge \
    -p '{"spec": {"devicePlugin": {"config": {"name": "time-slicing-config-all", "default": "any"}}}}'
```
