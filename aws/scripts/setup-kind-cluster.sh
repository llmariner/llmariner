#! /usr/bin/env bash
set -xe

# Boot kind cluster
cat <<EOF > kind-cluster-template.yaml
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
  labels:
    ingress-ready: true
  extraPortMappings:
  - containerPort: 80
    hostPort: 80
    protocol: TCP
  - containerPort: 443
    hostPort: 443
    protocol: TCP
{{- range \$gpu := until numGPUs }}
- role: worker
  extraMounts:
    - hostPath: /dev/null
      containerPath: /var/run/nvidia-container-devices/{{ \$gpu }}
{{- end }}
EOF
cat <<EOF > kind-cluster-values.yaml
numGPUs: 1
EOF
nvkind cluster list|grep 'No kind clusters found.' && { nvkind cluster create --config-template ./kind-cluster-template.yaml --config-values ./kind-cluster-values.yaml; sleep 5; }

# Add GPU operator
helm repo add nvidia https://helm.ngc.nvidia.com/nvidia
helm repo update
helm upgrade --install --wait \
     --namespace nvidia \
     --create-namespace \
     gpu-operator nvidia/gpu-operator \
     --set cdi.enabled=true \
     --set driver.enabled=false \
     --set toolkit.enabled=false

# Add Prometheus
cat <<EOF > prom-scrape-configs.yaml
- job_name: nvidia-dcgm
  scrape_interval: 5s
  static_configs:
  - targets: ['nvidia-dcgm-exporter.nvidia.svc:9400']
EOF
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update
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
