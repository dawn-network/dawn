package app

import (
	"fmt"
	"encoding/hex"
	"log"
	"github.com/tendermint/abci/types"
	"github.com/tendermint/go-merkle"
	dbm "github.com/tendermint/go-db"
)

type GlogChainApp struct {
	Db   		dbm.DB
	State 		merkle.Tree
	Height 		uint64
	TxCount   	uint64
}

type GlogVars struct {
	GlogApp *GlogChainApp
}

// global var
var GlogGlobal = GlogVars {}

func NewGlogChainApp() *GlogChainApp {
	log.Println("NewGlogChainApp")

	db := dbm.NewDB("state", dbm.GoLevelDBBackendStr, ".")

	lastBlock := LoadLastBlock(db)

	state := merkle.NewIAVLTree(0, db)
	state.Load(lastBlock.AppHash)

	log.Println("Loaded state", "block", lastBlock.Height, "root", state.Hash())

	return &GlogChainApp{ Db: db, State: state, Height: lastBlock.Height, TxCount: lastBlock.TxCount}
}

func (app *GlogChainApp) Info() (resInfo types.ResponseInfo) {
	log.Println("GlogChainApp.Info")

	resInfo = types.ResponseInfo{}
	resInfo.Data = "GlogChainApp"

	/////////////
	lastBlock := LoadLastBlock(app.Db)

	resInfo.LastBlockHeight = lastBlock.Height
	resInfo.LastBlockAppHash = lastBlock.AppHash
	return resInfo
}

func (app *GlogChainApp) SetOption(key string, value string) (logstr string) {
	log.Println("GlogChainApp.SetOption", key, value)
	return Exec_SetOption(key, value)
}

func (app *GlogChainApp) DeliverTx(tx []byte) types.Result {
	log.Println("GlogChainApp.DeliverTx")
	return Exec_DeliverTx(tx)
}

func (app *GlogChainApp) CheckTx(tx []byte) types.Result {
	log.Println("GlogChainApp.CheckTx")
	return Exec_CheckTx(tx)
}

func (app *GlogChainApp) Commit() types.Result {
	log.Println("GlogChainApp.Commit")

	appHash := app.State.Save()
	log.Println("Saved state", "root", hex.EncodeToString(appHash))

	lastBlock := LastBlockInfo {
		Height:  app.Height,
		AppHash: appHash, // this hash will be in the next block header
		TxCount: app.TxCount,
	}

	/////////////////////////////////////////
	SaveLastBlock(app.Db, lastBlock)
	log.Println("Saving block", "height", lastBlock.Height, "root", hex.EncodeToString(lastBlock.AppHash))

	return types.NewResultOK(appHash, "")
}

func (app *GlogChainApp) Query(query []byte) types.Result {
	log.Println("GlogChainApp.Query")
	return types.NewResultOK(nil, fmt.Sprintf("Query is not supported"))
}

func (app *GlogChainApp) InitChain(vals []*types.Validator) {
	log.Println("GlogChainApp.InitChain")
}

func (app *GlogChainApp) BeginBlock(hash []byte, header *types.Header) {
	log.Println("GlogChainApp.BeginBlock", "height=", header.Height, "hash=", hex.EncodeToString(hash))
	app.Height = header.Height
}

func (app *GlogChainApp) EndBlock(height uint64) (resEndBlock types.ResponseEndBlock) {
	log.Println("GlogChainApp.EndBlock", "height=", height, "size=", app.State.Size())
	return
}

