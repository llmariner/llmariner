#!/usr/bin/env bash
#
# Deploy External DNS so that we can automatically set up DNS names for ELBs (e.g., api.cloudnatix.com).
#
# See https://github.com/bitnami/charts/tree/master/bitnami/external-dns.

set -euo pipefail

basedir=$(dirname "$0")

helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update

namespace=external-dns
release=external-dns

helm upgrade \
  --install \
  --create-namespace \
  -n "${namespace}" \
  "${release}" bitnami/external-dns \
  --set provider=aws \
  --set domainFilters="{dev.llmo.cloudnatix.com}" \
  --set aws.zoneType=public \
  --set aws.assumeRoleArn=arn:aws:iam::803339316953:role/LLMOperatorDevExternalDNSMonitoring
