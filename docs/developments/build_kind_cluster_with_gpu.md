{# Build a Kind Cluster with GPU

## Overview

This document follows [Kevin Klues' instruction](https://github.com/klueska/kind-with-gpus-examples) to set up a Kind cluster with GPU.

The GitHub issue for the feature request is tracked in [here](https://github.com/kubernetes-sigs/kind/issues/3164).

## Procedure

### Step 1. Create an EC2 instance with GPU

- Create a `g5.4xlarge` instance. Ubuntu is used for this example.
- Allocate a 100 GiB of gp3 root volume as AI/ML models would require a large amount of data.

### Step 2. Install Required Tools

SSH into the instance and run the following tools:

- go
- make
- docker
- kind
- kubectl
- helm
- nvidia-driver
- nvidia-container-toolkit

```bash
# Go
curl -LO https://go.dev/dl/go1.22.2.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.22.2.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
export PATH=$PATH:$(go env GOPATH)/bin

# Make
sudo apt install make

# Docker
sudo apt-get update
sudo apt-get install ca-certificates curl
sudo install -m 0755 -d /etc/apt/keyrings
sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
sudo chmod a+r /etc/apt/keyrings/docker.asc
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu \
  $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | \
  sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
sudo apt-get update
sudo apt-get install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
sudo groupadd docker
sudo usermod -aG docker $USER
newgrp docker

# Kind
go install sigs.k8s.io/kind@v0.22.0

# Kubectl
curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
sudo install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl

# Helm
curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3
chmod 700 get_helm.sh
./get_helm.sh

# Nvidia driver
sudo apt install -y ubuntu-drivers-common
sudo ubuntu-drivers install

# Nvidia container toolkit
curl -fsSL https://nvidia.github.io/libnvidia-container/gpgkey | sudo gpg --dearmor -o /usr/share/keyrings/nvidia-container-toolkit-keyring.gpg \
  && curl -s -L https://nvidia.github.io/libnvidia-container/stable/deb/nvidia-container-toolkit.list | \
    sed 's#deb https://#deb [signed-by=/usr/share/keyrings/nvidia-container-toolkit-keyring.gpg] https://#g' | \
    sudo tee /etc/apt/sources.list.d/nvidia-container-toolkit.list
sudo apt-get update
sudo apt-get install -y nvidia-container-toolkit
```

### Step 3. Set up `nvkind`

```bash
sudo nvidia-ctk runtime configure --runtime=docker --set-as-default --cdi.enabled
sudo nvidia-ctk config --set accept-nvidia-visible-devices-as-volume-mounts=true --in-place
sudo systemctl restart docker
```


### Step 4. Create a Test Cluster

```bash
git clone https://github.com/klueska/kind-with-gpus-examples.git
cd kind-with-gpus-examples.git
make
./nvkind cluster create
./nvkind cluster list
./nvkind cluster print-gpus
```

```bash
helm repo add nvdp https://nvidia.github.io/k8s-device-plugin
helm repo update
helm upgrade -i \
  --namespace nvidia \
  --create-namespace \
  nvidia-device-plugin nvdp/nvidia-device-plugin
```

Create a test pod.
```
cat << EOF | kubectl apply -f -
apiVersion: v1
kind: Pod
metadata:
  name: gpu-test
spec:
  restartPolicy: OnFailure
  containers:
  - name: ctr
    image: ubuntu:22.04
    command: ["nvidia-smi", "-L"]
    resources:
      limits:
        nvidia.com/gpu: 1
EOF
```
