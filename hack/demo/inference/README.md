# Inference Autoscaling

## Preparation

Run on each EKS(auto-mode) worker cluster:

```bash
kubectl create namespace cloudnatix
kubectl apply -f gpu-nodepool.yaml

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

regKey=$(llma admin clusters register <CLUSTER>| sed -n 's/.*Registration Key: "\([^"]*\)".*/\1/p'
kubectl create secret -n cloudnatix generic cluster-registration-key --from-literal=regKey=${regKey}

helm upgrade \
  --install \
  -n cloudnatix \
  llmariner \
  oci://public.ecr.aws/cloudnatix/llmariner-charts/llmariner \
  -f ./llmariner-values.yaml
```

## Send Inference Requests

```bash
cat questions |\
    xargs -P 50 -n1 \
      llma chat completions create \
        --model lmstudio-community-phi-4-GGUF-phi-4-Q4_K_M.gguf \
        --role system \
        --completion
```

## Set up monitoring

Install DCGM Exporter:

```
kubectl create namespace nvidia

helm upgrade \
 --install \
 -n nvidia \
 dcgm-exporter \
 gpu-helm-charts/dcgm-exporter \
 --set serviceMonitor.enabled=false

kubectl apply -f ./dcgm-exporter-service.yaml
```

We might be able to use GPU Operator, but decided not to as it might
conflict with the existing setup in EKS Auto mode.

Then follow https://github.com/llmariner/llmariner/blob/main/provision/aws/scripts/setup-kind-cluster.sh#L44 to
install Prometheus and Grafana.

Please note that the Prometheus scraping config there is for a single
node. You need to specify individual nodes in the target.
