# Fine-tuning

## Install LLMariner

```bash
mkdir -p ~/.config/llmariner
cat << EOF > ~/.config/llmariner/config.yaml
version: v1
endpointUrl: https://api.llm.staging.cloudnatix.com/v1
auth:
  clientId: 0oa17m60zdJLsJUG14x7
  clientSecret: ""
  redirectUri: http://localhost:8084/callback
  issuerUrl: https://login.cloudnatix.com/oauth2/aus202ft6fhz9alff4x7
enableOkta: true
context:
  organizationId: org-z_PNhaYEjl1S6bWGh2RppPcy
  projectId: proj_UTxizYdNMTyDEh6tBNJ1SnJk
EOF

# Login with demo+gpu@cloudnatix.com
llma auth login
```

```bash
kubectl create namespace cloudnatix

export AWS_ACCESS_KEY_ID=...
export AWS_SECRET_ACCESS_KEY=...
kubectl create secret generic \
  aws \
  -n cloudnatix \
  --from-literal=accessKeyId=${AWS_ACCESS_KEY_ID} \
  --from-literal=secretAccessKey=${AWS_SECRET_ACCESS_KEY}

regKey=$(llma admin clusters register fine-tuning-demo | sed -n 's/.*Registration Key: "\([^"]*\)".*/\1/p')
kubectl create secret -n cloudnatix generic cluster-registration-key --from-literal=regKey=${regKey}
```

```bash
helm upgrade \
  --install \
  -n cloudnatix \
  llmariner \
  oci://public.ecr.aws/cloudnatix/llmariner-charts/llmariner \
  -f ./llmariner-values.yaml
```

## Create a training dataset

Follow <https://medium.com/@alexandros_chariton/how-to-fine-tune-llama-3-2-instruct-on-your-own-data-a-detailed-guide-e5f522f397d7>.

```bash
pip install datasets
python ./set_up_training_data.py

# Verify there is no format error
python ./validate_training_data_format.py

aws s3 cp training.jsonl s3://cloudnatix-installation-demo/training-data/training.jsonl

llma storage files create-link --object-path s3://cloudnatix-installation-demo/training-data/training.jsonl --purpose fine-tune
```

## Submit a fine-tuning job

```bash
export LLMARINER_API_KEY=...
python ./submit_fine_tuning_job.py <file ID>
```

## Use a generated model

Run chat completion with the generated model.
