# Running LLMariner on an EKS cluster with Karpenter

This summarizes how to run LLMariner on an EKS cluster that uses
Karpenter.

First, follow create an EKS cluster and install Karpenter by following

https://karpenter.sh/v0.37/getting-started/getting-started-with-karpenter/

Then run the following commands to create a `EC2NodeClass` and `NodePool`.

```bash
export GPU_AMI_ID="$(aws ssm get-parameter --name /aws/service/eks/optimized-ami/${K8S_VERSION}/amazon-linux-2-gpu/recommended/image_id --query Parameter.Value --output text)"

cat <<EOF | envsubst | kubectl apply -f -
apiVersion: karpenter.k8s.aws/v1beta1
kind: EC2NodeClass
metadata:
  name: default
spec:
  amiFamily: AL2 # Amazon Linux 2
  role: "KarpenterNodeRole-${CLUSTER_NAME}" # replace with your cluster name
  subnetSelectorTerms:
  - tags:
      karpenter.sh/discovery: "${CLUSTER_NAME}" # replace with your cluster name
  securityGroupSelectorTerms:
  - tags:
      karpenter.sh/discovery: "${CLUSTER_NAME}" # replace with your cluster name
  amiSelectorTerms:
   - id: "${GPU_AMI_ID}"

  blockDeviceMappings:
    - deviceName: /dev/xvda
      ebs:
        volumeSize: 1000Gi
        volumeType: gp3
        encrypted: true
        deleteOnTermination: true
EOF

cat <<EOF | envsubst | kubectl apply -f -
apiVersion: karpenter.sh/v1beta1
kind: NodePool
metadata:
  name: default
spec:
  template:
    spec:
      requirements:
      - key: "karpenter.k8s.aws/instance-family"
        operator: In
        values: ["m5"]
      - key: kubernetes.io/arch
        operator: In
        values: ["amd64"]
      - key: kubernetes.io/os
        operator: In
        values: ["linux"]
      nodeClassRef:
        apiVersion: karpenter.k8s.aws/v1beta1
        kind: EC2NodeClass
        name: default
  limits:
    cpu: 10000
  disruption:
    consolidationPolicy: WhenUnderutilized
    expireAfter: 720h # 30 * 24h = 720h
EOF
```

Please note the following points:

- GPU AMI ID needs to be included in `amiSelectorTerms`
- `blockDeviceMappings` is required to have a sufficient disk space for models

Karpenter uses this node group when a pod requests `nvidia.com/gpu`. I
don't know how Karpenter knows GPU instances are needed, but it might
have some internal logic.

When the node is initially created, the node doesn't have the
`nvidia.com/gpu` resource, but the GPU device plugins soon runs on the
node and make `nvidia.com/gpu` resource available to
`inference-manager-engine`.
