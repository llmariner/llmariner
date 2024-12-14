# Run LLMariner on AWS Inferentia

This is a note for running LLMariner on AWS Inferentia (`inf2` instances). This mostly follows
[this AWS blog post](https://aws.amazon.com/blogs/machine-learning/deploy-meta-llama-3-1-8b-on-aws-inferentia-using-amazon-eks-and-vllm/).

## Step 1. Provision `inf2` instances

Use Karpenter to provision `inf2` instances. Please note that
nodes need an AMI that is different from the case for Nvidia GPU.

```
export CLUSTER_NAME=kenji-karpenter-demo
export K8S_VERSION=1.31

export ACCELERATED_AMI_ID="$(aws ssm get-parameter --name /aws/service/eks/optimized-ami/${K8S_VERSION}/amazon-linux-2-gpu/recommended/image_id --query "Parameter.Value" --output text)"

cat << EOF | envsubst | kubectl apply -f -
apiVersion: karpenter.sh/v1
kind: NodePool
metadata:
  name: inferentia
spec:
  template:
    spec:
      requirements:
      - key: karpenter.sh/capacity-type
        operator: In
        values: ["on-demand"]
      nodeClassRef:
        group: karpenter.k8s.aws
        kind: EC2NodeClass
        name: inferentia
      requirements:
      - key: node.kubernetes.io/instance-type
        operator: In
        values:
        - inf2.24xlarge
      expireAfter: 720h
  disruption:
    consolidationPolicy: WhenEmptyOrUnderutilized
    consolidateAfter: 1m
---
apiVersion: karpenter.k8s.aws/v1
kind: EC2NodeClass
metadata:
  name: inferentia
spec:
  amiFamily: AL2
  role: "KarpenterNodeRole-${CLUSTER_NAME}"
  subnetSelectorTerms:
  - tags:
      karpenter.sh/discovery: "${CLUSTER_NAME}"
  securityGroupSelectorTerms:
  - tags:
      karpenter.sh/discovery: "${CLUSTER_NAME}"
  amiSelectorTerms:
  - id: "${ACCELERATED_AMI_ID}"
  blockDeviceMappings:
  - deviceName: /dev/xvda
    ebs:
      deleteOnTermination: true
      encrypted: true
      volumeSize: 256Gi
      volumeType: gp3
EOF
```

## Step 2. Install the Neuron device plugin and scheduling extension

Follow https://awsdocs-neuron.readthedocs-hosted.com/en/latest/containers/tutorials/k8s-setup.html#tutorial-k8s-env-setup-for-neuron.

```bash
helm upgrade --install neuron-helm-chart oci://public.ecr.aws/neuron/neuron-helm-chart \
  --set "npd.enabled=false"

helm upgrade --install neuron-helm-chart oci://public.ecr.aws/neuron/neuron-helm-chart \
  --set "scheduler.enabled=true" \
  --set "npd.enabled=false"
```

## Step 3. Build a VLLM container

Build a docker image for vLLM that runs on Inferentia nodes.

```bash
cat > Dockerfile <<EOF
FROM public.ecr.aws/neuron/pytorch-inference-neuronx:2.1.2-neuronx-py310-sdk2.20.0-ubuntu20.04

# Clone the vllm repository
RUN git clone https://github.com/vllm-project/vllm.git -b v0.6.4.post1

# Set the working directory
WORKDIR /vllm
# Set the environment variable
ENV VLLM_TARGET_DEVICE=neuron

# Install the dependencies
RUN python3 -m pip install -U -r requirements-neuron.txt
RUN python3 -m pip install .

# Modify the arg_utils.py file to support larger block_size option
RUN sed -i "/parser.add_argument('--block-size',/ {N;N;N;N;N;s/\[8, 16, 32\]/[8, 16, 32, 128, 256, 512, 1024, 2048, 4096, 8192]/}" vllm/engine/arg_utils.py

# Install ray
RUN python3 -m pip install ray
RUN pip install -U  triton>=3.0.0

# Set the entry point
ENTRYPOINT ["python3", "-m", "vllm.entrypoints.openai.api_server"]
EOF

docker build -t public.ecr.aws/cloudnatix/llmariner/vllm-openai:v0.6.4.post1-neuron .
docker push public.ecr.aws/cloudnatix/llmariner/vllm-openai:v0.6.4.post1-neuron
```

See also https://docs.vllm.ai/en/latest/getting_started/neuron-installation.html and
https://github.com/vllm-project/vllm/blob/main/Dockerfile.neuron.

## Step 4. Deploy LLMariner

Make the following configuration change to use the above container image for vLLM and
to allocate `aws.amazon.com/neuroncore` to model runtime pods.

```yaml
inference-manager-engine:
  runtime:
    runtimeImages:
      vllm: public.ecr.aws/cloudnatix/llmariner/vllm-openai:v0.6.2-neuron
  model:
    default:
      runtimeName: vllm
    overrides:
      meta-llama/Meta-Llama-3.3-70B-Instruct-fp8-dynamic:
       preloaded: true
        contextLength: 16384
        resources:
          limits:
            aws.amazon.com/neuroncore: 12
```

Note that vLLM does not currently work with `meta-llama/Meta-Llama-3.3-70B-Instruct-fp8-dynamic` due to the following error:

```python
ValueError: compressed-tensors quantization is currently not supported in Neuron Backend.
```
