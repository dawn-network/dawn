package service

import (
	"log"
	"io/ioutil"
	"github.com/baabeetaa/glogchain/app"
	"net/http"
)

// https://github.com/tendermint/tendermint/wiki/RPC
//Available endpoints:
//http://localhost:46657/dump_consensus_state
//http://localhost:46657/genesis
//http://localhost:46657/net_info
//http://localhost:46657/num_unconfirmed_txs
//http://localhost:46657/status
//http://localhost:46657/unconfirmed_txs
//http://localhost:46657/unsafe_stop_cpu_profiler
//http://localhost:46657/validators
//
//Endpoints that require arguments:
//http://localhost:46657/block?height=_
//http://localhost:46657/blockchain?minHeight=_&maxHeight=_
//http://localhost:46657/broadcast_tx_async?tx=_
//http://localhost:46657/broadcast_tx_sync?tx=_
//http://localhost:46657/dial_seeds?seeds=_
//http://localhost:46657/subscribe?event=_
//http://localhost:46657/unsafe_set_config?type=_&key=_&value=_
//http://localhost:46657/unsafe_start_cpu_profiler?filename=_
//http://localhost:46657/unsafe_write_heap_profile?filename=_
//http://localhost:46657/unsubscribe?event=_


// http://localhost:46657/status
func TmRpc_Status() (str_json_response string, err error) {
	var url_request string = app.GlogchainConfigGlobal.TmRpcLaddr + "/status"

	log.Println("url_request", url_request)

	resp, err := http.Get(url_request)
	if (err != nil) {
		log.Println(err.Error())
		return;
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if (err != nil) {
		log.Println(err.Error())
		return;
	}

	str_json_response = string(body[:])
	log.Println("json_response_string", str_json_response)

	return
}

func TmRpc_NetInfo() (str_json_response string, err error) {
	var url_request string = app.GlogchainConfigGlobal.TmRpcLaddr + "/net_info"

	log.Println("url_request", url_request)

	resp, err := http.Get(url_request)
	if (err != nil) {
		log.Println(err.Error())
		return;
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if (err != nil) {
		log.Println(err.Error())
		return;
	}

	str_json_response = string(body[:])
	log.Println("json_response_string", str_json_response)

	return
}

