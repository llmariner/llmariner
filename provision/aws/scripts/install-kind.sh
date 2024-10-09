#! /usr/bin/env bash
set -xe

: ${KIND_VERSION:=latest}

curl -Lo ./kind https://kind.sigs.k8s.io/dl/${KIND_VERSION}/kind-linux-amd64
sudo install -o root -g root -m 0755 kind /usr/local/bin/kind
rm ./kind
