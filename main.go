package main

import (
	"flag"

	. "github.com/tendermint/go-common"
	"github.com/tendermint/tmsp/server"
	"github.com/baabeetaa/glogchain/web"
	"github.com/baabeetaa/glogchain/config"
)

func main() {
	addrPtr := flag.String("addr", config.GlogchainConfigGlobal.TmspAddr, "Listen address")

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


func init() {
	config.ReadConfig()
}