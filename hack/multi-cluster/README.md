# Multi Cluster Deployment Testing

This directory contains scripts and config files for testing multi-cluster deployment.

To deploy:

```bash
./create_clusters.sh
./deploy.sh
```

The worker cluster uses an ExternalName service to reach the control plane.
Please note that the current service definition is for Mac/Windows (Docker Desktop).
See https://github.com/kubernetes-sigs/kind/issues/1200#issuecomment-130485579.

Please note that the endpoint address is http://localhost/v1, not http://localhost:8080/v1.
