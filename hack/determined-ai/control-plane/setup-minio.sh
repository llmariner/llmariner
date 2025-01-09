#!/usr/bin/env bash

set -euo pipefail

# Bucket name and the dummy region.
export S3_BUCKET_NAME=determined-ai

# Credentials for accessing the S3 bucket.
export AWS_ACCESS_KEY_ID=determined-ai-key
export AWS_SECRET_ACCESS_KEY=determined-ai-secret

helm upgrade --install --wait \
  --namespace minio \
  --create-namespace \
  minio oci://registry-1.docker.io/bitnamicharts/minio \
  --set auth.rootUser=minioadmin \
  --set auth.rootPassword=minioadmin \
  --set defaultBuckets="${S3_BUCKET_NAME}"

kubectl port-forward -n minio service/minio 9001 &

# Wait until the port-forwarding connection is established.
sleep 5

# Obtain the cookie and store in cookies.txt.
curl \
  http://localhost:9001/api/v1/login \
  --cookie-jar cookies.txt \
  --request POST \
  --header 'Content-Type: application/json' \
  --data @- << EOF
{
  "accessKey": "minioadmin",
  "secretKey": "minioadmin"
}
EOF

# Create a new API key.
curl \
  http://localhost:9001/api/v1/service-account-credentials \
  --cookie cookies.txt \
  --request POST \
  --header "Content-Type: application/json" \
  --data @- << EOF >/dev/null
{
  "name": "LLMariner",
  "accessKey": "$AWS_ACCESS_KEY_ID",
  "secretKey": "$AWS_SECRET_ACCESS_KEY",
  "description": "",
  "comment": "",
  "policy": "",
  "expiry": null
}
EOF

rm cookies.txt

kill %1
