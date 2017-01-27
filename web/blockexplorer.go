package web

import (
	"net/http"
	"github.com/baabeetaa/glogchain/service"
	"encoding/json"
	"log"
)

// Handerling /blockexplorer/* requests

/**
{
   "jsonrpc":"2.0",
   "id":"",
   "result":[
      32,
      {
         "node_info":{
            "pub_key":"557AA062B272B7F741440CEC354C5F7E53184D0CDF56FE03FD9EC3E2A8DDACBA",
            "moniker":"anonymous",
            "network":"dawn-chain",
            "remote_addr":"",
            "listen_addr":"10.0.0.11:46656",
            "version":"0.8.0",
            "other":[
               "wire_version=0.6.0",
               "p2p_version=0.3.5",
               "consensus_version=v1/0.2.2",
               "rpc_version=0.6.0/3",
               "rpc_addr=tcp://10.0.0.11:46657"
            ]
         },
         "pub_key":[
            1,
            "9C75D93028A6B52BA802CBC27A4EEB31E84917FB1B61E4C99B4F3091B9E3DCB9"
         ],
         "latest_block_hash":"C291B2DA5F7280CF371DA748349C177323CC0658",
         "latest_app_hash":"84201BB94504493163B3CDC002EA2C11F604A1BE",
         "latest_block_height":1714,
         "latest_block_time":1485509017681000000
      }
   ],
   "error":""
}
 */
func BlockExplorer_Status_Handler(w http.ResponseWriter, req *http.Request) {
	var data = map[string]interface{}{ }

	str_json_response, err := service.TmRpc_Status()
	if (err != nil) {
		render(w, "blockexplorer_status", ActionResult{Status: "error", Message: err.Error(), Data: data })
		return
	}

	err = json.Unmarshal([]byte(str_json_response), &data)
	if (err != nil) {
		log.Println(err.Error())
		return
	}

	data["json_str"] = str_json_response

	render(w, "blockexplorer_status", ActionResult{Status: "success", Message: "ok", Data: data })
	return
}

func RecentBlocksHandler(w http.ResponseWriter, req *http.Request) {

}

