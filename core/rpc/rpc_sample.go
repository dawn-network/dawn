package rpc

import (
	"log"
	"net/http"
)

type RpcSample int

type RpcSampleArgs struct {
	A, B int
}

type RpcSampleResult int


func (t *RpcSample) Multiply(r *http.Request, args *RpcSampleArgs, result *RpcSampleResult) error {
	log.Printf("Multiplying %d with %d\n", args.A, args.B)
	*result = RpcSampleResult(args.A * args.B)
	return nil
}
