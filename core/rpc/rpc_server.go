package rpc

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
	"net/http"
)


/**
https://haisum.github.io/2015/10/13/rpc-jsonrpc-gorilla-example-in-golang/


Test:
Send to http://10.0.0.11:8088/rpc
{
	"id": 1,
	"method": "RpcSample.Multiply",
	"params": [{"A": 5, "B": 2}]
}

Result;
{
	"result": 10,
	"error": null,
	"id": 1
}
 */
func StartRpcServer() {
	s := rpc.NewServer()

	s.RegisterCodec(json.NewCodec(), "application/json")
	s.RegisterCodec(json.NewCodec(), "application/json;charset=UTF-8")

	s.RegisterService(new(RpcSample), "")

	r := mux.NewRouter()
	r.Handle("/rpc", s)

	http.ListenAndServe(":8088", r)
}