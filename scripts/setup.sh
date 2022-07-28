#!/bin/bash
set -x

# Install necessary dependencies
sudo apt-get update -y
sudo DEBIAN_FRONTEND=noninteractive apt-get -y -o Dpkg::Options::="--force-confdef" -o Dpkg::Options::="--force-confold" dist-upgrade
sudo apt-get update
sudo apt-get -y -qq install curl wget git vim apt-transport-https ca-certificates
sudo add-apt-repository ppa:longsleep/golang-backports -y
sudo snap install go --classic 
#TODO: make postgres start on startup
# []  create a systemd unit file for app - starts on boot instead of doing it manually 
# []  Do I have to create a systemD file on terraform ssh server: this would be all through instruqt sandbox environment? or Can I do it through here ? 
# [] It would make sense to do it through here right instead of using a systemD file bc from here it will already be baked into our packer image 


go --version 

# Setup sudo to allow no-password sudo for "hashicorp" group and adding "terraform" user
sudo groupadd -r hashicorp
sudo useradd -m -s /bin/bash terraform
sudo usermod -a -G hashicorp terraform
sudo cp /etc/sudoers /etc/sudoers.orig
echo "terraform  ALL=(ALL) NOPASSWD:ALL" | sudo tee /etc/sudoers.d/terraform

# Installing SSH key
sudo mkdir -p /home/terraform/.ssh
sudo chmod 700 /home/terraform/.ssh
sudo cp /tmp/tf-packer.pub /home/terraform/.ssh/authorized_keys
sudo chmod 600 /home/terraform/.ssh/authorized_keys
sudo chown -R terraform /home/terraform/.ssh
sudo usermod --shell /bin/bash terraform

# Create GOPATH for Terraform user & download the webapp from github
#

sudo -H -i -u terraform -- env bash << EOF
whoami
echo ~terraform

cd /home/terraform

export GOPATH=/home/terraform/go
export PATH=$PATH:$GOROOT/bin:$GOPATH/bin
go install github.com/go-sql-driver/mysql
go install github.com/sabinlehaci/go-web-app@latest

EOF

