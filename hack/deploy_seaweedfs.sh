#!/usr/bin/env bash

set -euo pipefail
trap 'kill $(jobs -p)' EXIT

basedir=$(dirname "$0")

export AWS_ACCESS_KEY_ID=llmariner-key
export AWS_SECRET_ACCESS_KEY=llmariner-secret

# see https://github.com/seaweedfs/seaweedfs/wiki/Amazon-S3-API#public-access-with-anonymous-download for details.
cat <<EOF > s3-config.json
{
  "identities": [
    {
      "name": "me",
      "credentials": [
        {
          "accessKey": "$AWS_ACCESS_KEY_ID",
          "secretKey": "$AWS_SECRET_ACCESS_KEY"
        }
      ],
      "actions": [
        "Admin",
        "Read",
        "ReadAcp",
        "List",
        "Tagging",
        "Write",
        "WriteAcp"
      ]
    }
  ]
}
EOF

kubectl create namespace seaweedfs

# Create secrets.
kubectl create secret generic -n seaweedfs seaweedfs --from-file=s3-config.json
kubectl create secret generic -n llmariner aws \
  --from-literal=accessKeyId=${AWS_ACCESS_KEY_ID} \
  --from-literal=secretAccessKey=${AWS_SECRET_ACCESS_KEY}

rm s3-config.json

# deploy seaweedfs
kubectl apply -n seaweedfs -f "${basedir}"/seaweedfs.yaml
kubectl wait --timeout=60s --for=condition=ready pod -n seaweedfs -l app=seaweedfs

kubectl port-forward -n seaweedfs service/seaweedfs 8333 &
sleep 5

# Create a new bucket.
bucket_name=llmariner
aws --endpoint-url http://localhost:8333 s3 mb s3://${bucket_name}
