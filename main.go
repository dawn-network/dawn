package main

import (
	"flag"

	. "github.com/tendermint/go-common"
	"github.com/tendermint/tmsp/server"
	"glogchain/web"
)

func main() {
	addrPtr := flag.String("addr", "tcp://0.0.0.0:46658", "Listen address")

	flag.Parse()
	app := NewGlogChainApp()

	// Start the listener
	_, err := server.NewServer(*addrPtr, "grpc", app)
	if err != nil {
		Exit(err.Error())
	}


	// start web server on port 8000
	go web.StartWebServer()

	// Wait forever
	TrapSignal(func() {
		// Cleanup
	})
}
