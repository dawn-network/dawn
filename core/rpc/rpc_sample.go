package rpc

import (
	"log"
	"net/http"
)

type RpcSample struct {

}



/////////////////////
type RpcSample_Multiply_Args struct {
	A, B int
}

type RpcSample_Multiply_Result int

func (t *RpcSample) Multiply(r *http.Request, args *RpcSample_Multiply_Args, result *RpcSample_Multiply_Result) error {
	log.Println("Multiplying %d with %d", args.A, args.B)
	*result = RpcSample_Multiply_Result(args.A * args.B)
	return nil
}

/////////////////////
type RpcSample_Getty_Args struct {

}

type RpcSample_Getty_Result struct {
	A, B int
}

func (t *RpcSample) Getty(r *http.Request, args *RpcSample_Getty_Args, result *RpcSample_Getty_Result) error {
	log.Println("Getty")
	*result = RpcSample_Getty_Result {
		A: 111,
		B: 222,
	}
	return nil
}
