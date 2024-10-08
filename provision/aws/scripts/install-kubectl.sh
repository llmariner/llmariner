#! /usr/bin/env bash
set -xe

: ${KUBECTL_VERSION:=latest}

if [ "$KUBECTL_VERSION" = "latest" ]; then
    curl -Lo ./kubectl "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
else
    curl -Lo ./kubectl https://dl.k8s.io/release/${KUBECTL_VERSION}/bin/linux/amd64/kubectl
fi

sudo install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl
rm ./kubectl
