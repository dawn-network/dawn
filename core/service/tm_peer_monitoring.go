package service

import (
	"encoding/json"
	"time"
	"io/ioutil"
	"log"
	"github.com/dawn-network/glogchain/core/app"
	"net/http"
	"strings"
	"net/url"
	"net"
)

var Chain_Node_List []string


//////////////////////////////////////////////
// chain_peer_monitoring
func init() {
	log.Println("chain_peer_monitoring - init..")
	go chainpeer_monitoring()
}

func chainpeer_monitoring()  {
	for { // I will run forever
		//log.Println("chain_peer_monitoring")
		time.Sleep(time.Minute * 5)

		travelled_nodes := map[string]string{} 	// holding all nodes travelled


		// find peers on this node
		err := chainpeer_findpeers (app.GlogchainConfigGlobal.TmRpcLaddr, travelled_nodes, true)
		if (err != nil) {
			log.Println(err.Error())
			continue
		}

		// update to the list
		Chain_Node_List = []string{}
		for k := range travelled_nodes {
			u, err := url.Parse(k)
			if err != nil {
				log.Println(err.Error())
				break
			}

			host, _, err := net.SplitHostPort(u.Host)
			if err != nil {
				log.Println(err.Error())
				break
			}

			Chain_Node_List = append(Chain_Node_List, host)
		}
		log.Println(Chain_Node_List)
	}
}


func chainpeer_findpeers(node_address string, travelled_nodes map[string]string, with_nested bool) (err error) {
	// do nothing if this node is travelled
	if _, ok := travelled_nodes[node_address]; ok {
		return
	}
	travelled_nodes[node_address] = node_address // mard this node travelled

	/////////////////////////
	var url_request string = node_address + "/net_info"
	//log.Println("url_request", url_request)

	resp, err := http.Get(url_request)
	if (err != nil) {
		//log.Println(err.Error())
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if (err != nil) {
		//log.Println(err.Error())
		return
	}

	str_json_response := string(body[:])
	//log.Println("json_response_string", str_json_response)

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
			  "pub_key":"6E2ED97B8379D38D566279CCEDA0334CC23C57606648DF7142CF2DA743CE8DC1",
			  "moniker":"router",
			  "network":"dawn-test-two",
			  "remote_addr":"96.9.90.100:46656",
			  "listen_addr":"10.0.3.1:46656",
			  "version":"0.8.0",
			  "other":[
			     "wire_version=0.6.0",
			     "p2p_version=0.3.5",
			     "consensus_version=v1/0.2.2",
			     "rpc_version=0.6.0/3",
			     "rpc_addr=tcp://96.9.90.100:46657"
			  ]
		       },
		       "is_outbound":true,
		       "connection_status":{
			  "SendMonitor":{
			     "Active":true,
			     "Start":"2017-02-06T12:14:32.460Z",
			     "Duration":4944880000000,
			     "Idle":2780000000,
			     "Bytes":1501555,
			     "Samples":9731,
			     "InstRate":10,
			     "CurRate":20,
			     "AvgRate":304,
			     "PeakRate":5310,
			     "BytesRem":0,
			     "TimeRem":0,
			     "Progress":0
			  },
			  "RecvMonitor":{
			     "Active":true,
			     "Start":"2017-02-06T12:14:32.460Z",
			     "Duration":4944880000000,
			     "Idle":2780000000,
			     "Bytes":22122640,
			     "Samples":12698,
			     "InstRate":5,
			     "CurRate":77,
			     "AvgRate":4474,
			     "PeakRate":516450,
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
				"RecentlySent":0
			     },
			     {
				"ID":32,
				"SendQueueCapacity":100,
				"SendQueueSize":0,
				"Priority":5,
				"RecentlySent":822
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
				"RecentlySent":1551
			     },
			     {
				"ID":35,
				"SendQueueCapacity":2,
				"SendQueueSize":0,
				"Priority":1,
				"RecentlySent":56
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

	var res map[string]interface{}
	err = json.Unmarshal([]byte(str_json_response), &res)
	if (err != nil) {
		//log.Println(err.Error())
		return
	}

	res_result := res["result"].([]interface {})
	res_result_map := res_result[1].(map[string]interface {})
	res_peers := res_result_map["peers"].([]interface {})

	for _, v_peer := range res_peers {
		res_peer := v_peer.(map[string]interface {})
		node_info := res_peer["node_info"].(map[string]interface {})
		remote_addr := node_info["remote_addr"].(string)

		rpc_add_ip := strings.Split(remote_addr, ":")[0]
		//log.Println("rpc_add_ip", rpc_add_ip)


		peer_other := node_info["other"].([]interface {})
		for _, v_other_param := range peer_other {
			param := v_other_param.(string)
			if (strings.HasPrefix(param, "rpc_addr")) {
				rpc_add_port := strings.Split(param, ":")[2]
				//log.Println("rpc_add_port", rpc_add_port)
				rpc_add := "http://" + rpc_add_ip + ":" + rpc_add_port
				//log.Println("rpc_add", rpc_add)

				if (with_nested) {
					err = chainpeer_findpeers (rpc_add, travelled_nodes, true)
					if (err != nil) {
						//log.Println(err.Error())
						continue
					}
				} else {
					travelled_nodes[rpc_add] = rpc_add // mard this node travelled
				}
			}
		}
	}

	return
}
