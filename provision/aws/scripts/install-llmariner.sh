#! /usr/bin/env bash
set -xe

helmfile apply \
  --skip-diff-on-install \
  -l app!=fake-gpu-operator \
  --file ./dev/helmfile.yaml \
  --state-values-set llmariner.chartPath=../llmariner
