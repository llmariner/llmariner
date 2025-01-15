#!/usr/bin/env bash

set -euo pipefail

basedir=$(dirname "$0")

# Follow https://docs.konghq.com/kubernetes-ingress-controller/latest/get-started/
#
# The gateway API needs to be installed before Kong intallation as the Kong's helm chart behaves differently based on the presence of the gateway API
# (e.g., whether the cluster role includes HTTPRoutes).
#
# Use the experimental channel to install TCPRoute
kubectl apply -f https://github.com/kubernetes-sigs/gateway-api/releases/download/v1.2.0/experimental-install.yaml

helm repo add kong https://charts.konghq.com
helm repo update
helm upgrade --install kong kong/ingress  -f ./kong_values.yaml

kubectl apply -f "${basedir}"/gateway.yaml


# Enabling the feature gate is required to use TCPRoute.
# See https://docs.konghq.com/kubernetes-ingress-controller/latest/reference/feature-gates/.
kubectl patch deploy kong-controller --patch '{
  "spec": {
    "template": {
      "spec": {
        "containers": [
          {
            "name": "ingress-controller",
            "env": [
              {
                "name": "CONTROLLER_FEATURE_GATES",
                "value": "GatewayAlpha=true"
              }
            ]
          }
        ]
      }
    }
  }
}'

# Enable TCP request proxying.
#
# See https://docs.konghq.com/kubernetes-ingress-controller/latest/guides/services/tcp/#:~:text=Create%20TCP%20routing%20configuration%20for,Service%20that's%20running%20inside%20Kubernetes.
#
# TODO(kenji): Determined AI picks up port 50000 by default as we configure the port as the lowest available port.
# This works for this example demo setting, but in a real-world scenario, we need to support a range of ports.
kubectl patch deploy kong-gateway --patch '{
  "spec": {
    "template": {
      "spec": {
        "containers": [
          {
            "name": "proxy",
            "env": [
              {
                "name": "KONG_STREAM_LISTEN",
                "value": "0.0.0.0:50000"
              }
            ],
            "ports": [
              {
                "containerPort": 50000,
                "name": "det-50000",
                "protocol": "TCP"
              }
            ]
          }
        ]
      }
    }
  }
}'

kubectl patch service kong-gateway-proxy --patch '{
 "spec": {
   "ports": [
     {
       "name": "det-50000",
       "port": 50000,
       "protocol": "TCP",
       "targetPort": 50000
     }
   ]
 }
}'
