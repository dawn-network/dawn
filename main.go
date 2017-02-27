package main

import (
	"flag"

	. "github.com/tendermint/go-common"
	"github.com/tendermint/abci/server"
	"github.com/dawn-network/glogchain/web"
	. "github.com/dawn-network/glogchain/core/app"
	"log"
	cfg "github.com/tendermint/go-config"
	"github.com/tendermint/go-logger"
	"github.com/tendermint/tendermint/node"
	tmcfg "github.com/tendermint/tendermint/config/tendermint"
	"time"
	"os"
)

func main() {
	addrPtr := flag.String("addr", GlogchainConfigGlobal.TmspAddr, "Listen address")
	flag.Parse()

	GlogGlobal.GlogApp = NewGlogChainApp()

	if (GlogGlobal.GlogApp.Height == 0) { // genesis block
		GlogGlobal.GlogApp.SetOption("genesis.block/create.account", "pool/449F6F39391BC9E918CE51DB10F6FAADF65077263E715B17652425CD7827C814/1000000")
		//glogChainApp.SetOption("genesis.block/token.transfer", "17CE71F68874213CF40A512B162CBB3945EC35C9/1000000")

		GlogGlobal.GlogApp.SetOption("genesis.block/create.account", "jan/CDD6774218138DF657C7B0494894BBA596EB2ECCCC4C946D5ACEF3B5FCD2CE7B/1000")
		//glogChainApp.SetOption("genesis.block/token.transfer", "05D1D4B70CDA63A1A93FA381593A339BA9C67F19/1000")

		GlogGlobal.GlogApp.SetOption("genesis.block/create.account", "jake/488B8FF58E8E9868823C3388BAAB9C1F7CFCB3D7482376E7495639A1EC0F7407/1000")
		//glogChainApp.SetOption("genesis.block/token.transfer", "8DC49746AAB3E9A7D8546D1BF8497479B4A484CB/1000")

		GlogGlobal.GlogApp.SetOption("genesis.block/create.account", "tuan/EB3B42091EF6C2F8FA951319940C003BEC7AAE2336BD2AFABD6FB59EB4A3EF6E/1000")
		//glogChainApp.SetOption("genesis.block/token.transfer", "E405128ABE228A095EFF8D4940C66EC7A40AAA91/1000")
	}

	/////////////////////////////////////////////
	// Start the listener
	s, err := server.NewServer(*addrPtr, "grpc", GlogGlobal.GlogApp)
	if err != nil {
		Exit(err.Error())
	}

	/////////////////////////////////////////////
	// start web server on port 8000
	go web.StartWebServer()


	/////////////////////////////////////////////
	// start embedded tendermint
	go startTendermintNode()

	// Wait forever
	TrapSignal(func() {
		time.Sleep(3 * time.Second) // wait 3s for TM stopping

		// Cleanup
		s.Stop()
	})
}



var tm_config cfg.Config

/**
 Start Tendermint service as embedded mode.
 - Simpler for deploying
 - Hopefully avoid the panic bug when stop Glogchain before TM.
 */
func startTendermintNode()  {
	// Get configuration
	tm_config = tmcfg.GetConfig(GlogchainConfigGlobal.TmRoot)
	//parseFlags(config, args[1:]) // Command line overrides

	// set the log level
	logger.SetLogLevel(tm_config.GetString("log_level"))

	// wait sometime to make sure glogchain is up
	log.Println("Wait 10s to lauch Tendermint...")
	time.Sleep(time.Second * 10)
	node.RunNode(tm_config)
}

func init() {
	// to change the flags on the default logger
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// create tmp folder if need
	_, err := os.Stat("./tmp")
	if os.IsNotExist(err) {
		os.Mkdir("./tmp", os.ModePerm)
	}

	ReadConfig()
}