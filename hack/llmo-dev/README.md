# Hosting Configuration

This directory contains the configurations and scripts for deploying
LLM Operator and configure endpoint with https://api.dev.llmo.cloudnatix.com/v1.

```console
./deploy_cert_manager.sh
./deploy_external_dns.sh

./deploy_llm_operator.sh

kubectl apply -f kong_plugin.yaml
```

Cert Manager and External DNS require IAM role `arn:aws:iam::730335229895:role/LLMOperatorVMRole`.

# Limitation: External DNS

External DNS does not work in a Kind cluster in an EC2 instance. Even if we set
up [Cloud Provider Kind](https://github.com/kubernetes-sigs/cloud-provider-kind)
to be able to create a load balancer service, the IP of the load balancer is a private IP of the EC2 instances.

Hence we currently manually edit Route53 for setting up a DNS record for `api.dev.api.llmo.cloudnatix.com`.
