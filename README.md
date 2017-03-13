# glogchain: Let's upgrade the web

[![Build Status](http://163.172.170.63/api/badges/dawn-network/glogchain/status.svg)](http://163.172.170.63/dawn-network/glogchain)

Glogchain is dawn's first application specific blockchain.  Our blockchain stores a ledger and webtorrent/bittorrent hashes.  It will serve [webtorent](webtorrent.io) files to a single page web application that handles a number of different content types.  The back end API is provided by our network of validators.  Users upload content in text, audio, and video file formats and are able to share this content with their friends.  

You can download linux-only binaries here:  
https://github.com/dawn-network/glogchain/releases

Join our discourse chat at: https://discord.gg/8dWYbFS

<details>
<summary>Opinionated build for debian</summary>
```
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

```
</details>

<details>
<summary>test network</summary>
You can currently e-mail or send a google hangouts message to Jacob Gadikian at faddat@gmail.com for help getting onto one of our testnets as a validator, non-validator or light client.
</details>

<details>
<summary>Torrent Hashing: Intentionally Awesome</summary>

By storing hashes to our blockchain, and checking file hashes for validity, we are able to provide an immutability mechansim beyond our blockchain for the files that users upload.  

[![Router6d7376.md.png](http://www.steemimg.com/images/2017/02/07/Router6d7376.md.png)](http://www.steemimg.com/image/GhYv7)
</details>

<details>
<summary>Developer Experience</summary>

We have made a router, the Dawn R1, which happens to handily double as a computer which is equipped with a modern x86 CPU and adequate RAM and SSD storage.  This router makes an ideal development setup and comes pre-stocked with an opinionated golang development environment.  Developers do not need this router to participate, however fresh developers and experts alike will appreciate its isolated development environment that allows for fast, known-good development against our stack.  If you've any questions about the router or would like to buy one (sold at cost to developers who have made code commits to our projects) please contact Jacob Gadikian at faddat@gmail.com on google hangouts.
</details>

<details>
<summary>Privacy</summary>
![Screenshotfrom2017-02-0714-13-47dd71e.png](http://www.steemimg.com/images/2017/02/07/Screenshotfrom2017-02-0714-13-47dd71e.png)
This.

Privacy is implemented as follows:

* Public - Shared far and wide
* Private - Restricted to a key-holding group of individuals.  Users who do not possess the needed key are not allowed to decrypt private content.  We never possess the keys needed to unlock private content.
</details>

<details>
<summary>Coin and inflation</summary>
To ensure its survival in perpetuity, we have implemented a cryptocurrency system called Ray in glogchain.  One unit of currency is created with each block, forever.  This means that while early years will have a high inflation rate, actual currency supply after the 10th year or so will ahve relatively low inflation.  Given that content distribution and storage in this manner is at an infant state, we feel that ensuring a high enough validator count and a high (90%) rate of payments to creators based on hit count + eyeball-time will distribute the network's creative rays in a manner that reflects reality.  
</details>

<details>
<summary>Tendermint and Cosmos</summary>
We are of course huge fans of the [tendermint blockchain toolkit](github.com/tendermint/tendermint), and of the inter-network of blocckahins called [cosmos](github.com/tendermint/cosmos).  For more information, please see their whitepaper.  
</details>

<details>
<summary>Piracy</summary>
We encourage users to upload works to which they own the copyright.  Our seeders must unfortunately stop seeding files determined to contain copyrighted content not owned by the user.  Users may also choose to copyleft their content, or license it as they see fit.  The difference is that in our implementation, users drive decisionmaking about copyright, not a cabal of governmnet backed companies that have been around as long as recorded music.
</details>
