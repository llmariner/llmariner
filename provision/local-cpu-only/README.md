# Provisioner for a local CPU-only Kind cluster

This directory contains scripts that creates a Kind cluster locally and deploy LLMariner and other necessary components.

We use [Fake GPU operator](https://github.com/run-ai/fake-gpu-operator) so that LLMariner can run in a machine that doesn't have
any GPU. Ollama (not vLLM) is used as inference runtime.

Run:

```bash
./create_cluster.sh
./deploy.sh
```
