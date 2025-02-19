# Installation with CloudNatix

## Preparation

Create a S3 bucket:

```bash
aws s3 mb s3://cloudnatix-installation-demo --region us-west-2
```

Install Metric Server:

```bash
kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml
```

Install Nvidia Operator:

```bash
helm repo add nvidia https://helm.ngc.nvidia.com/nvidia
helm repo update
helm upgrade --install --wait \
  --namespace nvidia \
  --create-namespace \
  gpu-operator nvidia/gpu-operator \
  --set cdi.enabled=true \
  --set driver.enabled=false \
  --set toolkit.enabled=false
```

Create a secret for S3 and HuggingFace.

The AWS credentials are stored in the "dev" vault of 1Password.

```bash
kubectl create namespace cloudnatix

export AWS_ACCESS_KEY_ID=...
export AWS_SECRET_ACCESS_KEY=...
kubectl create secret generic \
  aws \
  -n cloudnatix \
  --from-literal=accessKeyId=${AWS_ACCESS_KEY_ID} \
  --from-literal=secretAccessKey=${AWS_SECRET_ACCESS_KEY}

kubectl create secret generic \
  huggingface-key \
  -n cloudnatix \
  --from-literal=apiKey=${HUGGING_FACE_HUB_TOKEN}
```

## Installation

Run:

```bash
export CNATIX_FEATURE_FLAG_LLMARINER=true
export CNATIX_GC_DOMAIN=staging.cloudnatix.com
export KUBECONFIG=<Vulter VKE kubeconfig>

# Login with demo+gpu@cloudnatix.com
llma auth login
llma admin clusters register my-demo-cluster

# Login with demo+gpu@cloudnatix.com
cnatix login
cnatix clusters configure
cnatix install
```

Select the `testing` channel until we release LLMariner to the `stable` channel.

If you just want to deploy LLMariner for testing, run:

```bash
registration_key=$(llma admin clusters register my-demo-cluster | sed -n 's/.*Registration Key: "\([^"]*\)".*/\1/p')
kubectl create secret -n cloudnatix generic cluster-registration-key --from-literal=regKey=${registration_key}
helm upgrade \
  --install \
  -n cloudnatix \
  llmariner \
  oci://public.ecr.aws/cloudnatix/llmariner-charts/llmariner \
  -f ./llmariner-values.yaml
```
