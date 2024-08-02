# Worker plane deployment

```bash
kubectl create namespace llm-operator

"${basedir}"/../../deploy_kong_internal.sh

export AWS_ACCESS_KEY_ID=llm-operator-key
export AWS_SECRET_ACCESS_KEY=llm-operator-secret
kubectl create secret generic -n llm-operator aws \
  --from-literal=accessKeyId=${AWS_ACCESS_KEY_ID} \
  --from-literal=secretAccessKey=${AWS_SECRET_ACCESS_KEY}

# Create a cluster registration credential
REGISTRATION_KEY=$(llmo admin clusters register worker-cluster | sed -n 's/.*Registration Key: "\([^"]*\)".*/\1/p')
kubectl create secret generic \
  -n llm-operator \
  cluster-registration-key \
  --from-literal=regKey="${REGISTRATION_KEY}"

./deploy_llm_operator.sh
```
