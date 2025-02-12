#! /usr/bin/env bash

set -xe

basedir=$(dirname "$0")

# Create a tenant  cluster
cat <<EOF | kind create cluster --name tenant-cluster --config -
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
featureGates:
  JobManagedBy: true
EOF

# Create two worker clusters.
cat <<EOF | kind create cluster --name gpu-worker-cluster-large --config -
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
- role: worker
- role: worker
- role: worker
- role: worker
EOF

cat <<EOF | kind create cluster --name gpu-worker-cluster-small --config -
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
- role: worker
EOF

# Set up the CLI configuration and login.
mkdir -p ~/.config/llmariner
cat << EOF > ~/.config/llmariner/config.yaml
version: v1
endpointUrl: https://api.llm.staging.cloudnatix.com/v1
auth:
  clientId: llmariner
  clientSecret: ZXhhbXBsZS1hcHAtc2VjcmV0
  redirectUri: http://127.0.0.1:5555/callback
  issuerUrl: https://api.llm.staging.cloudnatix.com/v1/dex
context:
  organizationId: org-z_PNhaYEjl1S6bWGh2RppPcy
  projectId: proj_UTxizYdNMTyDEh6tBNJ1SnJk
EOF

# Login with demp+gpu@cloudnatix.com
llma auth login

# Set up the worker clusters.
for cluster in "gpu-worker-cluster-large" "gpu-worker-cluster-small"; do
  kubectl config use-context "kind-${cluster}"

  # Deploy fake-gpu-operator
  kubectl label nodes --all --overwrite nodepool=default

  helm repo add fake-gpu-operator https://fake-gpu-operator.storage.googleapis.com
  helm repo update
  helm upgrade \
    --install \
    --create-namespace \
    -n nvidia \
    gpu-operator \
    fake-gpu-operator/fake-gpu-operator \
    --set topology.nodePools.default.gpuCount=8 \
    --set topology.nodePoolLabelKey=nodepool

  # Deploy LLMariner.
  registration_key=$(llma admin clusters register "${cluster}" | sed -n 's/.*Registration Key: "\([^"]*\)".*/\1/p')
  kubectl create namespace llmariner
  kubectl create secret -n llmariner generic cluster-registration-key --from-literal=regKey=${registration_key}
  helm upgrade \
    --install \
    -n llmariner \
    llmariner \
    oci://public.ecr.aws/cloudnatix/llmariner-charts/llmariner \
    -f "${basedir}"/llmarine-worker-cluster-values.yaml
done

# Set up the tenant cluster.
kubectl config use-context kind-tenant-cluster

tenant_api_key=$(llma auth api-keys create tenant -o 'Default Organization' --role tenant-system --service-account | sed -n 's/.*Secret: \(.*\)/\1/p')
kubectl create namespace llmariner
kubectl create secret -n llmariner generic syncer-api-key --from-literal=key=${tenant_api_key}

helm upgrade \
  --install \
  -n llmariner \
  llmariner \
  oci://public.ecr.aws/cloudnatix/llmariner-charts/llmariner \
  -f "${basedir}"/llmarine-tenant-cluster-values.yaml
