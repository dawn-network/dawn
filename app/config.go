package app

import (
	"fmt"
	"log"

	"path/filepath"
	"io/ioutil"
	"encoding/json"
)


type GlogCfg struct {
	//WebRootDir 		string	// dont need anymore because using https://github.com/elazarl/go-bindata-assetfs
	GlogchainWebAddr 	string
	TmspAddr 		string
	TmRpcLaddr 		string
	IpFsAPI			string
	IpFsGateway 		string
	TmRoot			string 	// where to store the tendermint data (chain + config for TM)
}

var GlogchainConfigGlobal = GlogCfg{}

func ReadConfig() error  {
	//config := GlogCfg{}

	filename, _ := filepath.Abs("./config.json")
	fmt.Printf("filename: %s\n", filename)

	yamlFile, err := ioutil.ReadFile(filename)

	if err != nil {
		panic(err)
		return err
	}

	err = json.Unmarshal(yamlFile, &GlogchainConfigGlobal)
	if err != nil {
		log.Fatalf("error: %v", err)
		return err
	}

	//fmt.Printf("config:\n%v\n\n", config)
	//fmt.Printf("Value: %s\n", config.Title)


	return nil
}
