set -xe

: ${HELMFILE_VERSION:=0.169.1}

curl -Lo ./helmfile.tar.gz "https://github.com/helmfile/helmfile/releases/download/v${HELMFILE_VERSION}/helmfile_${HELMFILE_VERSION}_linux_amd64.tar.gz"
tar -xvf helmfile.tar.gz helmfile
sudo install -o root -g root -m 0755 helmfile /usr/local/bin/helmfile
rm ./helmfile helmfile.tar.gz
