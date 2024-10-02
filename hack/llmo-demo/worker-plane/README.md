# Worker plane deployment

```bash
kubectl create namespace llmariner

export AWS_ACCESS_KEY_ID=llmariner-key
export AWS_SECRET_ACCESS_KEY=llmariner-secret
kubectl create secret generic -n llmariner aws \
  --from-literal=accessKeyId=${AWS_ACCESS_KEY_ID} \
  --from-literal=secretAccessKey=${AWS_SECRET_ACCESS_KEY}

# Create a cluster registration credential
REGISTRATION_KEY=$(llmo admin clusters register worker-cluster | sed -n 's/.*Registration Key: "\([^"]*\)".*/\1/p')
kubectl create secret generic \
  -n llmariner \
  cluster-registration-key \
  --from-literal=regKey="${REGISTRATION_KEY}"

./deploy_llmariner.sh
```
