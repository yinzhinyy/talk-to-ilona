#!/bin/sh

# prepare env
PROJECT_DIR=
if [ -n $1 ]; then
    PROJECT_DIR=~/go/src/$1
fi

rm -rf ${PROJECT_DIR}
go get -u github.com/yinzhinyy/talk-to-ilona
go install $1

sudo mkdir -p /etc/talk
sudo cp ${PROJECT_DIR}/configs/* /etc/talk/

sudo mkdir -p /usr/lib/systemd/system
sudo cp ${PROJECT_DIR}/init/talk.service /usr/lib/systemd/system/
sudo systemctl daemon-reload

sudo systemctl restart talk
