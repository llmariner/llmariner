#! /usr/bin/env bash
set -xe

helmfile apply \
  --skip-diff-on-install \
  -l app!=fake-gpu-operator \
  --file ./dev/helmfile.yaml.gotmpl \
  --state-values-set llmariner.chartPath=../llmariner
