package web

import (
	"net/http"
	"github.com/gorilla/mux"
	"log"
	"github.com/dawn-network/glogchain/core/service"
	"strconv"
)


/**
 Get the torrent file and manipulate announces and webseeds before return to client
 See more service.Create_Torrent_From_Local_File
 */
func TorrentFileHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	ipfshash, ok := vars["ipfshash"]
	if !ok {
		serve404(w)
		return
	}

	log.Println("ipfshash", ipfshash)
	data, err := service.Get_Torrent_From_Ipfs(ipfshash)
	if (err != nil) {
		serve400(w)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename="+ ipfshash + ".torrent")
	w.Header().Set("Content-Type", "application/x-bittorrent")
	w.Header().Set("Content-Length", strconv.FormatInt(int64(len(data)), 10)) //Get file size as a string
	w.Write(data)
}
