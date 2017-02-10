package web

import (
	"net/http"
	"github.com/dawn-network/glogchain/service"
	"encoding/json"
	"log"
	"strconv"
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

	///////
	//result := data["result"].([]interface{})
	//result_info := result[1].(map[string]interface{})
	//latest_block_height := result_info["latest_block_height"]
	//log.Println("latest_block_height", latest_block_height)

	data["json_str"] = str_json_response

	render(w, "blockexplorer_status", ActionResult{Status: "success", Message: "ok", Data: data })
	return
}

func BlockExplorer_NetInfo_Handler(w http.ResponseWriter, req *http.Request) {
	var data = map[string]interface{}{ }

	str_json_response, err := service.TmRpc_NetInfo()
	if (err != nil) {
		render(w, "blockexplorer_netinfo", ActionResult{Status: "error", Message: err.Error(), Data: data })
		return
	}

	err = json.Unmarshal([]byte(str_json_response), &data)
	if (err != nil) {
		log.Println(err.Error())
		return
	}

	data["json_str"] = str_json_response


	///////
	bChain_Node_List, err := json.Marshal(service.Chain_Node_List)
	if err != nil {
		render(w, "blockexplorer_netinfo", ActionResult{Status: "error", Message: err.Error(), Data: data })
		return;
	}
	data["Chain_Node_List"] = string(bChain_Node_List)

	render(w, "blockexplorer_netinfo", ActionResult{Status: "success", Message: "ok", Data: data })
	return
}

func BlockExplorer_ConsensusState_Handler(w http.ResponseWriter, req *http.Request) {
	var data = map[string]interface{}{ }

	str_json_response, err := service.TmRpc_Dump_Consensus_State()
	if (err != nil) {
		render(w, "blockexplorer_dump_consensus_state", ActionResult{Status: "error", Message: err.Error(), Data: data })
		return
	}

	err = json.Unmarshal([]byte(str_json_response), &data)
	if (err != nil) {
		log.Println(err.Error())
		return
	}

	data["json_str"] = str_json_response

	render(w, "blockexplorer_dump_consensus_state", ActionResult{Status: "success", Message: "ok", Data: data })
	return
}

/**
{
   "jsonrpc":"2.0",
   "id":"",
   "result":[
      3,
      {
         "block_meta":{
            "hash":"DB14E6A3C2E71B09B2D2E71A36E249D220D2FEB0",
            "header":{
               "chain_id":"dawn-chain",
               "height":3276,
               "time":"2017-01-27T16:00:56.829Z",
               "num_txs":0,
               "last_block_id":{
                  "hash":"000F140AF954A5D04B8D68486BFF0460D88E60B2",
                  "parts":{
                     "total":1,
                     "hash":"DD410CFF60A7081744D7D14339C1566D184FEBCD"
                  }
               },
               "last_commit_hash":"EC5A3209EB22532F12D909B33AB6F4F07F848024",
               "data_hash":"",
               "validators_hash":"FEB47FCB45234702C8ABB4C860C772D1B53A310C",
               "app_hash":"84201BB94504493163B3CDC002EA2C11F604A1BE"
            },
            "parts_header":{
               "total":1,
               "hash":"083CAFEE9509E2E90C2F73D380F58CEEB67293A1"
            }
         },
         "block":{
            "header":{
               "chain_id":"dawn-chain",
               "height":3276,
               "time":"2017-01-27T16:00:56.829Z",
               "num_txs":0,
               "last_block_id":{
                  "hash":"000F140AF954A5D04B8D68486BFF0460D88E60B2",
                  "parts":{
                     "total":1,
                     "hash":"DD410CFF60A7081744D7D14339C1566D184FEBCD"
                  }
               },
               "last_commit_hash":"EC5A3209EB22532F12D909B33AB6F4F07F848024",
               "data_hash":"",
               "validators_hash":"FEB47FCB45234702C8ABB4C860C772D1B53A310C",
               "app_hash":"84201BB94504493163B3CDC002EA2C11F604A1BE"
            },
            "data":{
               "txs":[

               ]
            },
            "last_commit":{
               "blockID":{
                  "hash":"000F140AF954A5D04B8D68486BFF0460D88E60B2",
                  "parts":{
                     "total":1,
                     "hash":"DD410CFF60A7081744D7D14339C1566D184FEBCD"
                  }
               },
               "precommits":[
                  {
                     "validator_address":"D0F50DE71B2EBA74290C2F3DBB2AEAA54B289C8A",
                     "validator_index":0,
                     "height":3275,
                     "round":0,
                     "type":2,
                     "block_id":{
                        "hash":"000F140AF954A5D04B8D68486BFF0460D88E60B2",
                        "parts":{
                           "total":1,
                           "hash":"DD410CFF60A7081744D7D14339C1566D184FEBCD"
                        }
                     },
                     "signature":[
                        1,
                        "62F6AB8748786C0229E38C4EB437485CB4F75C03F868A5896FB54616BBF295F8DEF0971E5C04EAAB21D1D395516553002FCC6F363423A81A84FA0DC954A57802"
                     ]
                  }
               ]
            }
         }
      }
   ],
   "error":""
}
 */
func BlockExplorer_Block_Handler(w http.ResponseWriter, req *http.Request) {
	var data = map[string]interface{}{ }

	pHeight := req.FormValue("height")

	var height int64
	var err error

	if (pHeight != "") {
		height, err = strconv.ParseInt(pHeight, 10, 64)
		if (err != nil) {
			render(w, "blockexplorer_block", ActionResult{Status: "error", Message: err.Error(), Data: data })
			return
		}
	} else {
		height = getLastBlockHeight()
	}

	str_json_response, err := service.TmRpc_Block(height)
	if (err != nil) {
		render(w, "blockexplorer_block", ActionResult{Status: "error", Message: err.Error(), Data: data })
		return
	}

	err = json.Unmarshal([]byte(str_json_response), &data)
	if (err != nil) {
		log.Println(err.Error())
		return
	}

	data["json_str"] = str_json_response

	render(w, "blockexplorer_block", ActionResult{Status: "success", Message: "ok", Data: data })
	return
}

//////////////////////////////////////////////
// local functions

func getLastBlockHeight() int64  {
	var data = map[string]interface{}{ }

	str_json_response, err := service.TmRpc_Status()
	if (err != nil) {
		log.Println(err.Error())
		return 0
	}

	err = json.Unmarshal([]byte(str_json_response), &data)
	if (err != nil) {
		log.Println(err.Error())
		return 0
	}

	/////
	result := data["result"].([]interface{})
	result_info := result[1].(map[string]interface{})
	latest_block_height := result_info["latest_block_height"]
	log.Println("latest_block_height", latest_block_height)

	return int64(latest_block_height.(float64))
}
