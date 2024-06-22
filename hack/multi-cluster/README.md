# Multi Cluster Deployment Testing

This directory contains scripts and config files for testing multi-cluster deployment.

To deploy:

```bash
./create_clusters.sh
./deploy.sh
```

We create a NodePort service and set up the external port mapping in Kind
so that the worker service of `session-manager-server`
can be reachable from the worker cluster.

The worker cluster uses an ExternalName service to reach the control plane.
Please note that the current service definition is for Mac/Windows (Docker Desktop).
See https://github.com/kubernetes-sigs/kind/issues/1200#issuecomment-130485579.
