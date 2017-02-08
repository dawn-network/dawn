package service

import (
	//"log"
	//"github.com/ipfs/go-ipfs-api"
	"io"
	//"github.com/baabeetaa/glogchain/app"
)

// NOTICE: dont use IPFS anymore, maybe replace by https://github.com/syncthing/syncthing

// https://ipfs.io/docs/api/#apiv0add
func Ipfs_add(r io.Reader) (mhash string, err error) {
	//log.Println("Ipfs_add")
	//
	//s := shell.NewShell(app.GlogchainConfigGlobal.IpFsAPI)
	//
	//mhash, err = s.Add(r)
	//
	//log.Println("mhash", mhash)

	return
}