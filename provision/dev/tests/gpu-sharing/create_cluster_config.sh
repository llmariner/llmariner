#!/usr/bin/env bash

set -euo pipefail

echo "Creating a cluster config..."
llma admin clusters config create Default --time-slicing-gpus 2

llma admin clusters config get Default
