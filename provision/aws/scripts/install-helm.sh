#! /usr/bin/env bash
set -xe

: ${HELM_VERSION:=latest}

curl https://baltocdn.com/helm/signing.asc | gpg --dearmor | sudo tee /usr/share/keyrings/helm.gpg > /dev/null
echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/helm.gpg] https://baltocdn.com/helm/stable/debian/ all main" | sudo tee /etc/apt/sources.list.d/helm-stable-debian.list
sudo apt update

if [ "$HELM_VERSION" = "latest" ]; then
    sudo apt install -y helm
else
    sudo apt install -y helm=${HELM_VERSION}
fi
