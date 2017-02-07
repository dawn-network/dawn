# glogchain: better communities

Glogchain is dawn's first application specific blockchain.  Our blockchain stores a ledger and webtorrent hashes.  It will serve [webtorent](webtorrent.io) files to a single page web application that handles a number of different content types.  The back end API is provided by our network of validators.  Users upload content in text, audio, and video file formats and are able to share this content with their friends.  

By building this repository you can use an incomplete version of our network. 

## test network
If you wish to join the test network as a validator, please contact Jacob Gadikian at faddat@gmail.com on Google Hangouts.  

## Torrent Hashing: Intentionally Awesome

By storing hashes to our blockchain, and checking file hashes for validity, we are able to provide an immutability mechansim beyond our blockchain for the files that users upload.  

## Development Kit

We have made a router, the Dawn R1, which happens to handily double as a computer which is equipped with a modern x86 CPU and adequate RAM and SSD storage.  This router makes an ideal development setup and comes pre-stocked with an opinionated golang development environment.  Developers do not need this router to participate, however fresh developers and experts alike will appreciate its isolated development environment.  If you've any questions about the router or would like to buy one (sold at cost to developers) please contact Jacob Gadikian at faddat@gmail.com on google hangouts.  

## Privacy
![Screenshotfrom2017-02-0714-13-47dd71e.png](http://www.steemimg.com/images/2017/02/07/Screenshotfrom2017-02-0714-13-47dd71e.png)
This.

Privacy is implemented as follows:

* Public - Shared far and wide
* Private - Restricted to a key-holding group of individuals.  Users who do not possess the needed key are not allowed to decrypt private content.  

## seed nodes
Please see @baabeetaa's [guide to creating non-validator nodes](https://github.com/baabeetaa/glogchain/wiki/Create-local-testnet).  Seed addresses are listed there and you should be able to join our test network.  

## binaries
Binaries for OSX, Linux, and windows will be available shortly.  

## Coin and inflation
To ensure its survival in perpetuity, we have implemented a cryptocurrency system called Ray in glogchain.  One unit of currency is created with each block, forever.  This means that while early years will have a high inflation rate, actual currency supply after the 10th year or so will ahve relatively low inflation.  Given that content distribution and storage in this manner is at an infant state, we feel that ensuring a high enough validator count and a high (90%) rate of payments to creators based on hit count + eyeball-time will distribute the network's creative rays in a manner that reflects reality.  

## Tendermint and Cosmos
We are of course huge fans of the tendermint blockchain toolkit, and of the inter-network of blocckahins called cosmos.  For more information, please see their whitepaper.  


