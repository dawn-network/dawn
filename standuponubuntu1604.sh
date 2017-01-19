apt update
apt install curl git mercurial make binutils bison gcc build-essential wget
apt upgrade -y
bash < <(curl -s -S -L https://raw.githubusercontent.com/moovweb/gvm/master/binscripts/gvm-installer)
source /root/.gvm/scripts/gvm
gvm install go1.8rc1 -B
gvm use go 1.8rc1 --default
go get github.com/baabeetaa/glogchain
cd $GOPATH/src/github.com/baabeetaa/glogchain
glogchain

