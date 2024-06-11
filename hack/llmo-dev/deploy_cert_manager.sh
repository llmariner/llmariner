#!/usr/bin/env bash
#
# Deploy Cert Manager (https://cert-manager.io/).
#
# See https://cert-manager.io/docs/installation/kubernetes/#steps for the install steps.

set -euo pipefail

basedir=$(dirname "$0")

kubectl apply --validate=false -f https://github.com/jetstack/cert-manager/releases/download/v1.12.0/cert-manager.crds.yaml

helm repo add jetstack https://charts.jetstack.io
helm repo update

release=cert-manager
namespace=cert-manager
helm upgrade \
  --install \
  --wait \
  --create-namespace \
  -n "${namespace}" \
  "${release}" jetstack/cert-manager --version v1.12.0

kubectl apply -f "${basedir}"/letsencrypt-clusterissuer-staging.yaml
kubectl apply -f "${basedir}"/letsencrypt-clusterissuer-prod.yaml
