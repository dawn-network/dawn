package service

import (
	"log"
	"github.com/ipfs/go-ipfs-api"
	"io"
	"github.com/baabeetaa/glogchain/app"
)

//const (
//	shellUrl     = "localhost:5001"
//)

// https://ipfs.io/docs/api/#apiv0add
func Ipfs_add(r io.Reader) (mhash string, err error) {
	log.Println("Ipfs_add")

	s := shell.NewShell(app.GlogchainConfigGlobal.IpFsAPI)

	//mhash, err = s.Add(bytes.NewBufferString("Hello IPFS Shell tests"))
	mhash, err = s.Add(r)

	log.Println("mhash", mhash)

	return
}