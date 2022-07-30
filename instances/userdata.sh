#! /usr/bin/env bash
set -xeuo pipefail


export TMDB = af2b4e4b7c2c224650dfad4faa2de6ff
export POSTGRES_URL= postgres://postgres:mysecretpassword@localhost:5432/postgres

apt-get -q -y install postgresql 

sudo apt-get install -y \
    apt-transport-https \
    ca-certificates \
    curl \
    software-properties-common 

curl -fsSL https://yum.dockerproject.org/gpg | sudo apt-key add - 

sudo add-apt-repository \
    "deb https://apt.dockerproject.org/repo/ \
    ubuntu-$(lsb_release -cs) \
    main" 

sudo apt-get update
sudo apt-get -y install docker-engine 

# add current user to docker group so there is no need to use sudo when running docker
sudo usermod -aG docker $(whoami)

docker run --rm -p 5432:5432 -e POSTGRES_PASSWORD=mysecretpassword -d postgres:14


# Start the web server
/opt/webapp/server
