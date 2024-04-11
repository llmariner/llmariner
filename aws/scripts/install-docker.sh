#! /usr/bin/env bash
set -xe

: ${DOCKER_VERSION:=latest}

# install packages
sudo install -m 0755 -d /etc/apt/keyrings
GPG_KEY=/etc/apt/keyrings/docker.gpg
[ -f ${GPG_KEY} ] || sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o ${GPG_KEY}
sudo chmod a+r $GPG_KEY
echo \
  "deb [arch="$(dpkg --print-architecture)" signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
  "$(. /etc/os-release && echo "$VERSION_CODENAME")" stable" | \
  sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
sudo apt update
if [ "$DOCKER_VERSION" = "latest" ]; then
  sudo apt install -y docker-ce docker-ce-cli containerd.io
else
  sudo apt install -y docker-ce=${DOCKER_VERSION} docker-ce-cli=${DOCKER_VERSION} containerd.io
fi

# setup docker group
sudo groupadd docker || :
sudo usermod -aG docker $USER

# start docker
sudo systemctl daemon-reload
sudo systemctl enable docker
sudo systemctl restart docker
