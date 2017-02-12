package service

import (
	"github.com/anacrolix/torrent/metainfo"
	"time"
	"log"
	"github.com/anacrolix/torrent/bencode"
	"os"
	"bufio"
	"github.com/dawn-network/glogchain/app"
	"net/http"
	"bytes"
)

var builtinAnnounceList = [][]string{
	{"udp://tracker.openbittorrent.com:80"},
	{"udp://tracker.internetwarriors.net:1337"},
	{"udp://tracker.leechers-paradise.org:6969"},
	{"udp://tracker.coppersurfer.tk:6969"},
	{"udp://exodus.desync.com:6969"},
	{"wss://tracker.btorrent.xyz"},
	{"wss://tracker.openwebtorrent.com"},
	{"wss://tracker.fastcast.nz"},
}

/**
 output is <filepath>.torrent

 Note that announces and webseeds can be changed at runtime. So that this torrent file should not be served directly,
 Instead, it need to be grapped and manipulate the torrent by adding announces and webseeds before return to clients.

 filepath: is the path to torrent file
 mhash_upload: is the hash of the file uploaded to ipfs
 */
func Create_Torrent_From_Local_File(filepath string, mhash_upload string) (err error)  {
	mi := metainfo.MetaInfo{
		AnnounceList: builtinAnnounceList,
	}

	mi.Comment = mhash_upload // store the ipfs hash to the Comment so that we know how to access from ipfs
	mi.CreatedBy = "github.com/dawn-network/glogchain"
	mi.CreationDate = time.Now().Unix()

	info := metainfo.Info {
		PieceLength: 512 * 1024,
	}
	err = info.BuildFromFilePath(filepath)

	if err != nil {
		log.Println(err.Error())
		return
	}
	mi.InfoBytes, err = bencode.Marshal(info)
	if err != nil {
		log.Println(err.Error())
		return
	}

	mi.URLList = []string { }

	var file_torrent = filepath + ".torrent"
	log.Println("file_torrent", file_torrent)
	f, err := os.Create(file_torrent)
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer f.Close()

	w := bufio.NewWriter(f)

	err = mi.Write(w)
	if err != nil {
		log.Println(err.Error())
		return
	}

	w.Flush()

	log.Println("mi", mi.HashInfoBytes().HexString())
	return
}

func Get_Torrent_From_Ipfs(mhash string) (data []byte, err error)  {
	/////////////////////////////////
	// get the file first
	url := app.GlogchainConfigGlobal.IpFsGateway + "/ipfs/" + mhash
	//log.Println("url", url)
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer resp.Body.Close()


	/////////////////////////////////
	// load to torrent
	//r := bytes.NewReader(data)
	mi, err := metainfo.Load(resp.Body)
	if err != nil {
		log.Println(err.Error())
		return
	}
	//info, err := mi.UnmarshalInfo()
	//if err != nil {
	//	log.Println(err.Error())
	//	return
	//}

	/////////////////////////////////
	// manipulate the torrent file
	mhash_upload := mi.Comment // see Create_Torrent_From_Local_File

	// TODO: need to update webseeds as list of current validator nodes
	webseed := app.GlogchainConfigGlobal.IpFsGateway + "/ipfs/" + mhash_upload
	mi.URLList = []string { webseed }

	// write torrent to raw data
	var buf bytes.Buffer
	w := bufio.NewWriter(&buf)
	err = mi.Write(w)
	if err != nil {
		log.Println(err)
		return
	}
	w.Flush()

	data = buf.Bytes()

	return
}