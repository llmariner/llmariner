#! /usr/bin/env bash
set -xe

# Nvidia driver
sudo apt install -y alsa-base ubuntu-drivers-common
sudo ubuntu-drivers autoinstall

# Nvidia container toolkit
GPG_KEY=/usr/share/keyrings/nvidia-container-toolkit-keyring.gpg
[ -f ${GPG_KEY} ] || curl -fsSL https://nvidia.github.io/libnvidia-container/gpgkey | sudo gpg --dearmor -o ${GPG_KEY} \
  && curl -s -L https://nvidia.github.io/libnvidia-container/stable/deb/nvidia-container-toolkit.list | \
    sed 's#deb https://#deb [signed-by=/usr/share/keyrings/nvidia-container-toolkit-keyring.gpg] https://#g' | \
    sudo tee /etc/apt/sources.list.d/nvidia-container-toolkit.list
sudo apt update
sudo apt install -y nvidia-container-toolkit

# Setup nvkind
sudo nvidia-ctk runtime configure --runtime=docker --set-as-default --cdi.enabled
sudo nvidia-ctk config --set accept-nvidia-visible-devices-as-volume-mounts=true --in-place
sudo systemctl restart docker

# Build nvkind
WORKSPACE=$HOME/kind-with-gpus-examples
[ -d $WORKSPACE ] || git clone https://github.com/klueska/kind-with-gpus-examples.git $WORKSPACE
cd $WORKSPACE
source $HOME/.profile
make
sudo install -o root -g root -m 0755 nvkind /usr/local/bin/nvkind
