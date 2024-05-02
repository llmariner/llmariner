# Deploy a k8s cluster with kubespray

## Requirements

- [Ansible](https://docs.ansible.com/ansible/latest/installation_guide/intro_installation.html)
- [virtualenv](https://virtualenv.pypa.io/en/latest/)

## Step 1. Clone kubespray repository

```bash
git clone https://github.com/kubernetes-sigs/kubespray.git
cd kubespray
git checkout 91dea023ae9ea829be5ae9458fb8b5799a48d779
```

## Step 2. Set up inventory

Copy inventory and variable files to the kubespray repository.

```bash
cp -R ../inventory/llmo inventory/
```

Then, edit the `inventory.ini` file according to your environment.

> [!NOTE]
> As an example, the llmo directory contains inventory file for a all-in-one single node cluster.
> See the [official document](https://kubespray.io/#/docs/ansible?id=inventory) for the inventory details.

## Step 3. Install required packages

```bash
python3 -m venv venv
source venv/bin/activate
pip install -U -r requirements.txt
pip install jsonschema
```

## Step 4. Deploy a k8s cluster

```bash
ansible-playbook -bvi inventory/llmo/inventory.ini \
    --user=`<USER (e.g., ubuntu)>` \
    --private-key=`<PATH/TO/YOUR/KEY>`
```

## Step 5. Access to the machine and set up components

```bash
# Copy k8s configuration file
mkdir .kube
sudo cp /etc/kubernetes/admin.conf .kube/config
sudo chown $USER:$USER .kube/config
kubectl cluster-info

# Add GPU operator
helm repo add nvidia https://helm.ngc.nvidia.com/nvidia
helm repo update
helm upgrade --install --wait \
     --namespace nvidia \
     --create-namespace \
     gpu-operator nvidia/gpu-operator \
     --set cdi.enabled=true

# Add Prometheus
cat <<EOF > prom-values.yaml
server:
  persistentVolume:
    enabled: false
prometheus-pushgateway:
  enabled: false
alertmanager:
  enabled: false
EOF
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
     -f prom-values.yaml \
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
```
