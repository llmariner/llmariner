#! /usr/bin/env bash

set -xe

basedir=$(dirname "$0")

for i in {1..20}; do
  cat <<EOF | kubectl apply --context kind-tenant-cluster -f -
  apiVersion: batch/v1
  kind: Job
  metadata:
    name: gpu-job-${i}
  spec:
    managedBy: cloudnatix.com/job-controller
    template:
      spec:
        containers:
        - name: gpu-job
          image: ubuntu
          command:
          - /bin/sleep
          - "300"
          resources:
            limits:
              nvidia.com/gpu: 1
        restartPolicy: Never
EOF
  sleep 1
done
