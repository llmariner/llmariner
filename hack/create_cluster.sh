#!/usr/bin/env bash

set -euo pipefail

cluster_name="llm-operator-demo"

kind create cluster --name "${cluster_name}"
