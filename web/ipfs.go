package web

import (
	"net/http"
	"github.com/dawn-network/glogchain/service"
	"log"
	"os"
	"io"
)

func TestIpFsHandler(w http.ResponseWriter, req *http.Request) {
	if (req.Method != "POST") {
		//crutime := time.Now().Unix()
		//h := md5.New()
		//io.WriteString(h, strconv.FormatInt(crutime, 10))
		//token := fmt.Sprintf("%x", h.Sum(nil))

		render(w, "test_ipfs", ActionResult{Status: "success", Message: "ok", Data: map[string]interface{}{ }})
		return
	}


	////////////////
	// handling POST here
	req.ParseMultipartForm(32 << 20)
	file_upload, handler, err := req.FormFile("uploadfile")
	if (err != nil) {
		log.Println(err.Error())
		return
	}
	defer file_upload.Close()

	log.Println(handler.Filename)

	// save to local temp file
	////fmt.Fprintf(w, "%v", handler.Header)
	filepath_relative := "./tmp/" + handler.Filename
	filepath_torrent := filepath_relative + ".torrent"


	file_temp, err := os.OpenFile(filepath_relative, os.O_WRONLY|os.O_CREATE, 0666)
	if (err != nil) {
		log.Println(err.Error())
		return
	}
	defer file_temp.Close()
	io.Copy(file_temp, file_upload)

	//raw, err := ioutil.ReadAll(file)
	//if (err != nil) {
	//	log.Println(err.Error())
	//	return;
	//}

	defer func() {
		//////////////////////////////////////////////
		// delete the temp file
		file_upload.Close()
		file_temp.Close()
		err = os.Remove(filepath_relative)
		err = os.Remove(filepath_torrent)
	}()


	//////////////////////////////////////////////
	// upload to IPFS network
	mhash_upload, err := service.Ipfs_add_file(filepath_relative)
	if (err != nil) {
		log.Println(err.Error())
		return
	}


	//////////////////////////////////////////////
	// generate torrent file and upload to ipfs
	err = service.Create_Torrent_From_Local_File(filepath_relative, mhash_upload)
	if (err != nil) {
		log.Println(err.Error())
		return
	}

	// upload torrent to IPFS network

	mhash_torrent, err := service.Ipfs_add_file(filepath_torrent)
	if (err != nil) {
		log.Println(err.Error())
		return
	}

	// use {{ Config_IpFsGateway }}
	// "http_gateway": app.GlogchainConfigGlobal.IpFsGateway + "/ipfs/" + mhash_file,
	render(w, "test_ipfs", ActionResult{Status: "success", Message: "ok",
		Data: map[string]interface{}{
			"mhash_upload": mhash_upload,
			"mhash_torrent": mhash_torrent,
		},
	})
}
