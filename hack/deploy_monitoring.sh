#!/usr/bin/env bash

set -euo pipefail

# Add Prometheus
cat <<EOF > prom-scrape-configs.yaml
- job_name: nvidia-dcgm
  scrape_interval: 5s
  static_configs:
  - targets: ['nvidia-dcgm-exporter.nvidia.svc:9400']
EOF
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update
# TODO(kenji): Consider deploy prometheus-operator instead.
helm upgrade --install --wait \
     --namespace monitoring \
     --create-namespace \
     --set-file extraScrapeConfigs=prom-scrape-configs.yaml \
     prometheus prometheus-community/prometheus

# Add Grafana with DCGM dashboard
cat <<EOF > grafana-values.yaml
datasources:
 datasources.yaml:
   apiVersion: 1
   datasources:
   - name: Prometheus
     type: prometheus
     url: http://prometheus-server
     isDefault: true
dashboardProviders:
  dashboardproviders.yaml:
    apiVersion: 1
    providers:
    - name: 'default'
      orgId: 1
      folder: 'default'
      type: file
      disableDeletion: true
      editable: true
      options:
        path: /var/lib/grafana/dashboards/standard
dashboards:
  default:
    nvidia-dcgm-exporter:
      gnetId: 12239
      datasource: Prometheus
EOF
helm repo add grafana https://grafana.github.io/helm-charts
helm repo update
helm upgrade --install --wait \
     --namespace monitoring \
     --create-namespace \
     -f grafana-values.yaml \
     grafana grafana/grafana
