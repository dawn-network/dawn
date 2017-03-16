#!/bin/bash
echo -n "Enter your public IP address and press [ENTER]: "
read $PUBIP
echo -n "node listening address $PUBIP:46656"
echo -n "HTTP address [ENTER]: $PUBIP:80"
apt-get update
apt-get -y upgrade
apt-get -y install  curl git mercurial make binutils bison gcc build-essential protobuf-compiler
curl -s -S -L https://raw.githubusercontent.com/moovweb/gvm/master/binscripts/gvm-installer| bash
source /root/.gvm/scripts/gvm
gvm install go1.8 -B -pb
gvm use go1.8 --default
mkdir -p $GOPATH/bin
go get -u github.com/Masterminds/glide
mkdir $GOPATH/src/github.com/tendermint
git clone https://github.com/tendermint/tendermint/ $GOPATH/src/github.com/tendermint/tendermint
cd $GOPATH/src/github.com/tendermint/tendermint
git branch develop
make install
git clone https://github.com/dawn-network/glogchain/ $GOPATH/src/github.com/dawn-network/glogchain
cd $GOPATH/src/github.com/dawn-network/glogchain
sed -ie 's/10.0.0.11/$PUBIP/g' config.json
git branch develop
glide install
go build .
go install .
go get github.com/ipfs/go-ipfs
cd $GOPATH/src/github.com/ipfs/go-ipfs
make install
cd ~/
ipfs init
ipfs cat /ipfs/QmVLDAhCY3X9P2uRudKAryuQFPM5zqA3Yij1dY8FpGbL7T/readme
tendermint init
echo "Please ensure that you have set up any needed forwarding.  This script detects your public IP address, but there are many good reasons why your machine may not be using a public IP address.  We've no love for NAT, so if you're running a validator, please, no NAT."
