apt install build-essential bison
bash < <(curl -s -S -L https://raw.githubusercontent.com/moovweb/gvm/master/binscripts/gvm-installer)
source /home/$USER/.gvm/scripts/gvm
gvm install go1.8 -B -pb
gvm use go1.8 --default
mkdir $GOPATH/bin
go get -u github.com/Masterminds/glide
go get -u github.com/tendermint/tendermint/...
cd $GOPATH/src/github.com/tendermint/tendermint
git branch develop
make install
go get -u github.com/dawn-network/glogchain/...
git branch develop
cd $GOPATH/src/github.com/dawn-network/glogchain
glide install
go build .
go install .
