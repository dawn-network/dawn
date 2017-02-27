#!/bin/sh

######################################
# generate embeded resources (assetfs)

#cd $GOPATH/src/github.com/dawn-network/glogchain/web
#cd web

go-bindata-assetfs -pkg web webcontent/...



######################################
# build

#cd ../
#go build
