# Example Setup for Gateway API

This directory contains an example setup for using Gateway API for routing requests from `session-manager-agent` to Jupyter Notebooks.

First build a Kind cluster and deploy the components as usual.

```bash
../create_cluster.sh
../deploy.sh
```

Then deploy a Kong in the `llm-operator` namespace and configure it as a gateway.

```bash
./deploy_kong_as_gateway.sh
```

Then deploy LLM Operator.

```bash
./deploy_llm_operator.sh
```

You can create a Jupyter Notebook and verify the access.

```bash
rm -rf ~/.config/llmo
llmo auth login
llmo workspace notebooks create my-nb
llmo workspace notebooks open my-nb
```
