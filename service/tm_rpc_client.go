package service

import (
	"log"
	"io/ioutil"
	//"github.com/dawn-network/glogchain/app"
	"net/http"
	"fmt"
	"github.com/dawn-network/glogchain/types"
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

//http://localhost:46657/dump_consensus_state
func TmRpc_Dump_Consensus_State() (str_json_response string, err error) {
	var url_request string = types.GlogchainConfigGlobal.TmRpcLaddr + "/dump_consensus_state"
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

// http://localhost:46657/status
func TmRpc_Status() (str_json_response string, err error) {
	/**
	{
	   "jsonrpc":"2.0",
	   "id":"",
	   "result":[
	      32,
	      {
		 "node_info":{
		    "pub_key":"D46A04C9526C959322143E1E7B0BAEF3239C6EE9B7D901847351DE4F1AD1C2F3",
		    "moniker":"baabeetaa",
		    "network":"dawn-test-two",
		    "remote_addr":"",
		    "listen_addr":"192.168.1.25:46656",
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
		    "F314F023F14E793867BE7B06FFAA2E1BBC0C91697AF57A3BC3C1DC4A26C94A6B"
		 ],
		 "latest_block_hash":"0FF370DC66330062F145FCE13C452FE2D32B2478",
		 "latest_app_hash":"C6B2A339CC3BE59F77498AC56A907A9A15F1B919",
		 "latest_block_height":9390,
		 "latest_block_time":1486362892925000000
	      }
	   ],
	   "error":""
	}
	 */

	var url_request string = types.GlogchainConfigGlobal.TmRpcLaddr + "/status"
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

//http://localhost:46657/net_info
func TmRpc_NetInfo() (str_json_response string, err error) {
	/**
	{
	   "jsonrpc":"2.0",
	   "id":"",
	   "result":[
	      33,
	      {
		 "listening":true,
		 "listeners":[
		    "Listener(@192.168.1.25:46656)"
		 ],
		 "peers":[
		    {
		       "node_info":{
			  "pub_key":"ABEB9E806D958D53E8B3B845A9A8BD8344628C963D17302C038E558F23871A1E",
			  "moniker":"baabeetaa",
			  "network":"dawn-test-two",
			  "remote_addr":"163.172.170.188:46656",
			  "listen_addr":"163.172.170.188:46656",
			  "version":"0.8.0",
			  "other":[
			     "wire_version=0.6.0",
			     "p2p_version=0.3.5",
			     "consensus_version=v1/0.2.2",
			     "rpc_version=0.6.0/3",
			     "rpc_addr=tcp://10.2.18.157:46657"
			  ]
		       },
		       "is_outbound":true,
		       "connection_status":{
			  "SendMonitor":{
			     "Active":true,
			     "Start":"2017-02-06T22:47:17.040Z",
			     "Duration":211560000000,
			     "Idle":80000000,
			     "Bytes":68697,
			     "Samples":651,
			     "InstRate":19,
			     "CurRate":66,
			     "AvgRate":325,
			     "PeakRate":6570,
			     "BytesRem":0,
			     "TimeRem":0,
			     "Progress":0
			  },
			  "RecvMonitor":{
			     "Active":true,
			     "Start":"2017-02-06T22:47:17.040Z",
			     "Duration":211600000000,
			     "Idle":40000000,
			     "Bytes":8052156,
			     "Samples":984,
			     "InstRate":5542,
			     "CurRate":10962,
			     "AvgRate":38054,
			     "PeakRate":517300,
			     "BytesRem":0,
			     "TimeRem":0,
			     "Progress":0
			  },
			  "Channels":[
			     {
				"ID":48,
				"SendQueueCapacity":1,
				"SendQueueSize":0,
				"Priority":5,
				"RecentlySent":0
			     },
			     {
				"ID":64,
				"SendQueueCapacity":100,
				"SendQueueSize":0,
				"Priority":5,
				"RecentlySent":2720
			     },
			     {
				"ID":32,
				"SendQueueCapacity":100,
				"SendQueueSize":0,
				"Priority":5,
				"RecentlySent":0
			     },
			     {
				"ID":33,
				"SendQueueCapacity":100,
				"SendQueueSize":0,
				"Priority":10,
				"RecentlySent":0
			     },
			     {
				"ID":34,
				"SendQueueCapacity":100,
				"SendQueueSize":0,
				"Priority":5,
				"RecentlySent":0
			     },
			     {
				"ID":35,
				"SendQueueCapacity":2,
				"SendQueueSize":0,
				"Priority":1,
				"RecentlySent":0
			     }
			  ]
		       }
		    },
		    {
		       "node_info":{
			  "pub_key":"649D2DE0F2EFB42592B1E6A496A1535EB18BDC452006EA2D264D6D2DFA0DEDBC",
			  "moniker":"baabeetaa",
			  "network":"dawn-test-two",
			  "remote_addr":"51.15.47.174:46656",
			  "listen_addr":"51.15.47.174:46656",
			  "version":"0.8.0",
			  "other":[
			     "wire_version=0.6.0",
			     "p2p_version=0.3.5",
			     "consensus_version=v1/0.2.2",
			     "rpc_version=0.6.0/3",
			     "rpc_addr=tcp://10.8.41.15:46657"
			  ]
		       },
		       "is_outbound":true,
		       "connection_status":{
			  "SendMonitor":{
			     "Active":true,
			     "Start":"2017-02-06T22:47:17.180Z",
			     "Duration":211500000000,
			     "Idle":940000000,
			     "Bytes":53775,
			     "Samples":346,
			     "InstRate":10,
			     "CurRate":25,
			     "AvgRate":254,
			     "PeakRate":6570,
			     "BytesRem":0,
			     "TimeRem":0,
			     "Progress":0
			  },
			  "RecvMonitor":{
			     "Active":true,
			     "Start":"2017-02-06T22:47:17.180Z",
			     "Duration":211500000000,
			     "Idle":400000000,
			     "Bytes":6304930,
			     "Samples":744,
			     "InstRate":2630,
			     "CurRate":4468,
			     "AvgRate":29811,
			     "PeakRate":517180,
			     "BytesRem":0,
			     "TimeRem":0,
			     "Progress":0
			  },
			  "Channels":[
			     {
				"ID":48,
				"SendQueueCapacity":1,
				"SendQueueSize":0,
				"Priority":5,
				"RecentlySent":0
			     },
			     {
				"ID":64,
				"SendQueueCapacity":100,
				"SendQueueSize":0,
				"Priority":5,
				"RecentlySent":1901
			     },
			     {
				"ID":32,
				"SendQueueCapacity":100,
				"SendQueueSize":0,
				"Priority":5,
				"RecentlySent":0
			     },
			     {
				"ID":33,
				"SendQueueCapacity":100,
				"SendQueueSize":0,
				"Priority":10,
				"RecentlySent":0
			     },
			     {
				"ID":34,
				"SendQueueCapacity":100,
				"SendQueueSize":0,
				"Priority":5,
				"RecentlySent":0
			     },
			     {
				"ID":35,
				"SendQueueCapacity":2,
				"SendQueueSize":0,
				"Priority":1,
				"RecentlySent":0
			     }
			  ]
		       }
		    }
		 ]
	      }
	   ],
	   "error":""
	}
	 */

	var url_request string = types.GlogchainConfigGlobal.TmRpcLaddr + "/net_info"
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

//http://localhost:46657/block?height=_
func TmRpc_Block(height int64) (str_json_response string, err error) {
	/**
	{
	   "jsonrpc":"2.0",
	   "id":"",
	   "result":[
	      3,
	      {
		 "block_meta":{
		    "hash":"30E0AF5A0705CC7588D62849E2886159AC463934",
		    "header":{
		       "chain_id":"dawn-test-two",
		       "height":1,
		       "time":"2017-02-05T23:14:27.249Z",
		       "num_txs":0,
		       "last_block_id":{
			  "hash":"",
			  "parts":{
			     "total":0,
			     "hash":""
			  }
		       },
		       "last_commit_hash":"",
		       "data_hash":"",
		       "validators_hash":"3C8E72B2FD6DC780703CB592B6476896352D7485",
		       "app_hash":""
		    },
		    "parts_header":{
		       "total":1,
		       "hash":"038CEAAAF346027CE8378445DA4A70338BE54F42"
		    }
		 },
		 "block":{
		    "header":{
		       "chain_id":"dawn-test-two",
		       "height":1,
		       "time":"2017-02-05T23:14:27.249Z",
		       "num_txs":0,
		       "last_block_id":{
			  "hash":"",
			  "parts":{
			     "total":0,
			     "hash":""
			  }
		       },
		       "last_commit_hash":"",
		       "data_hash":"",
		       "validators_hash":"3C8E72B2FD6DC780703CB592B6476896352D7485",
		       "app_hash":""
		    },
		    "data":{
		       "txs":[

		       ]
		    },
		    "last_commit":{
		       "blockID":{
			  "hash":"",
			  "parts":{
			     "total":0,
			     "hash":""
			  }
		       },
		       "precommits":[

		       ]
		    }
		 }
	      }
	   ],
	   "error":""
	}
	 */
	var url_request string = fmt.Sprintf("%s/block?height=%d", types.GlogchainConfigGlobal.TmRpcLaddr, height)
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
