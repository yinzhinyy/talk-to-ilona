#!/bin/sh

# prepare env
PROJECT_DIR=
if [ -n $1 ]; then
    PROJECT_DIR=~/go/src/$1
fi

sudo mkdir -p /etc/talk
sudo cp ${PROJECT_DIR}/configs/* /etc/talk/

sudo mkdir -p /usr/lib/systemd/system
sudo cp ${PROJECT_DIR}/init/talk.service /usr/lib/systemd/system/
sudo systemctl daemon-reload

go install $1

sudo systemctl restart talk
