#! /usr/bin/env bash
set -xeuo pipefail


export TMDB = af2b4e4b7c2c224650dfad4faa2de6ff
export POSTGRES_URL= postgres://postgres:mysecretpassword@localhost:5432/postgres

apt-get -q -y install postgresql 
sudo bash 
apt install -y curl ca-certificates gnupg 
curl https://www.postgresql.org/media/keys/ACCC4CF8.asc | sudo apt-key add -
sh -c 'echo "deb http://apt.postgresql.org/pub/repos/apt $(lsb_release -cs)-pgdg main" > /etc/apt/sources.list.d/postgresql.list' 
apt update -y
apt install -y postgresql-client-11
pg_basebackup -V

yes | sudo apt-get install docker.io 

# add current user to docker group so there is no need to use sudo when running docker
sudo usermod -aG docker $(whoami)

docker run --rm -p 5432:5432 -e POSTGRES_PASSWORD=mysecretpassword -d postgres:14


# Start the web server
/opt/webapp/server


