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
