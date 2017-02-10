package web

import (
	"net/http"
	"github.com/dawn-network/glogchain/service"
	"log"
	"github.com/dawn-network/glogchain/app"
	"io/ioutil"
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
	file, handler, err := req.FormFile("uploadfile")
	if (err != nil) {
		log.Println(err.Error())
		return
	}
	defer file.Close()

	log.Println(handler.Filename)

	////fmt.Fprintf(w, "%v", handler.Header)
	//filpath_relative := "./tmp/" + handler.Filename
	//f, err := os.OpenFile(filpath_relative, os.O_WRONLY|os.O_CREATE, 0666)
	//if (err != nil) {
	//	log.Println(err.Error())
	//	return
	//}
	//defer f.Close()
	//io.Copy(f, file)

	raw, err := ioutil.ReadAll(file)
	if (err != nil) {
		log.Println(err.Error())
		return;
	}


	// upload to IPFS network
	mhash, err := service.Ipfs_add(raw)
	if (err != nil) {
		log.Println(err.Error())
		return
	}

	render(w, "test_ipfs", ActionResult{Status: "success", Message: "ok",
		Data: map[string]interface{}{
			"mhash": mhash,
			"http_gateway": app.GlogchainConfigGlobal.IpFsGateway + "/ipfs/" + mhash,
		},
	})
}
