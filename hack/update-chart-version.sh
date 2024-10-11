#!/usr/bin/env bash

set -eu -o pipefail

ver=$1
name=llmariner

echo "Update chart and app version for $name" \
    && sed -i "s/^version:.*$/version: ${ver}/" deployments/$name/Chart.yaml \
    && sed -i "s/^appVersion:.*$/appVersion: ${ver}/" deployments/$name/Chart.yaml \
    && echo "=> done!"
