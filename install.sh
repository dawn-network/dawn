#!/bin/bash

#Ensure that we're root
if [[ $UID != 0 ]]; then
    echo "Please run this script with sudo:"
    echo "sudo $0 $*"
    exit 1
fi

#Change to /root because script is relative to that path
cd /root

#Update Machine and install dependencies
apt-get update
apt-get -y upgrade
apt-get -y install  curl git mercurial make binutils bison gcc build-essential protobuf-compiler

#install GVM and go1.8
curl -s -S -L https://raw.githubusercontent.com/moovweb/gvm/master/binscripts/gvm-installer | bash
source /root/.gvm/scripts/gvm
gvm install go1.8 -B -pb
gvm use go1.8 --default

#Install glide
mkdir -p $GOPATH/bin
go get -u github.com/Masterminds/glide

#Install Tendermint
mkdir $GOPATH/src/github.com/tendermint
git clone https://github.com/tendermint/tendermint/ $GOPATH/src/github.com/tendermint/tendermint
cd $GOPATH/src/github.com/tendermint/tendermint
git branch develop
make install

#Install glogchain
go get -u github.com/dawn-network/glogchain
cd $GOPATH/src/github.com/dawn-network/glogchain
mkdir ~/.glogchain
go build .
go install .

#Install IPFS
go get github.com/ipfs/go-ipfs
cd $GOPATH/src/github.com/ipfs/go-ipfs
make install

#Initialize IPFS and tendermint
cd ~/
ipfs init
ipfs cat /ipfs/QmVLDAhCY3X9P2uRudKAryuQFPM5zqA3Yij1dY8FpGbL7T/readme
tendermint init
echo "Please ensure that you have set up any needed forwarding.  This script detects your public IP address, but there are many good reasons why your machine may not be using a public IP address.  We've no love for NAT, so if you're running a validator, please, no NAT."
