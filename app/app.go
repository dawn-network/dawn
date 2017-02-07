package app

import (
	"fmt"
	"log"
	"github.com/tendermint/abci/types"
	"github.com/tendermint/go-merkle"
	dbm "github.com/tendermint/go-db"
	"encoding/hex"
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

//var current_BlockHeader *types.Header
var savedAppHash []byte = []byte("savedAppHash")
var hasTx bool		// indicate current block has any TX

func NewGlogChainApp() *GlogChainApp {
	db := dbm.NewDB("state", dbm.GoLevelDBBackendStr, ".")
	state := merkle.NewIAVLTree(0, db)

	lastBlock := LoadLastBlock(db)
	state.Load(lastBlock.AppHash)
	savedAppHash = state.Hash()

	log.Println("NewGlogChainApp", "block", lastBlock.Height, "root: ", hex.EncodeToString(state.Hash()))

	return &GlogChainApp{ Db: db, State: state, Height: lastBlock.Height, TxCount: lastBlock.TxCount}
}

func (app *GlogChainApp) Info() (resInfo types.ResponseInfo) {
	resInfo = types.ResponseInfo{}
	resInfo.Data = "GlogChainApp"

	lastBlock := LoadLastBlock(app.Db)
	resInfo.LastBlockHeight = lastBlock.Height
	resInfo.LastBlockAppHash = lastBlock.AppHash

	log.Println("Info", "block", lastBlock.Height, "root", hex.EncodeToString(lastBlock.AppHash))

	return resInfo
}

func (app *GlogChainApp) SetOption(key string, value string) (logstr string) {
	//log.Println("GlogChainApp.SetOption", key, value)
	hasTx = true
	return Exec_SetOption(key, value)
}

func (app *GlogChainApp) DeliverTx(tx []byte) types.Result {
	//log.Println("GlogChainApp.DeliverTx")
	hasTx = true
	return Exec_DeliverTx(tx)
}

func (app *GlogChainApp) CheckTx(tx []byte) types.Result {
	//log.Println("GlogChainApp.CheckTx")
	return Exec_CheckTx(tx)
}

func (app *GlogChainApp) Commit() types.Result {
	//log.Println("GlogChainApp.Commit")

	// dont need to save state if there is no transaction or state doesnt change
	if hasTx || (app.Height <= 1)  { // (app.Height % 10000 == 0
		savedAppHash = app.State.Save()
		lastBlock := LastBlockInfo {
			Height:  app.Height,
			AppHash: savedAppHash, 			// this hash will be in the next block header
			TxCount: app.TxCount,
		}

		SaveLastBlock(app.Db, lastBlock)
	}

	return types.NewResultOK(savedAppHash, "")
}

func (app *GlogChainApp) Query(query []byte) types.Result {
	log.Println("GlogChainApp.Query")
	return types.NewResultOK(nil, fmt.Sprintf("Query is not supported"))
}

func (app *GlogChainApp) InitChain(vals []*types.Validator) {
	log.Println("GlogChainApp.InitChain")
}

func (app *GlogChainApp) BeginBlock(hash []byte, header *types.Header) {
	//log.Println("GlogChainApp.BeginBlock", "height=", header.Height, "hash=", hex.EncodeToString(hash))
	app.Height = header.Height
	//CloneValue(header, &current_BlockHeader)
	//current_BlockHeader = header
	hasTx = false		// enter new block, reset value
}

func (app *GlogChainApp) EndBlock(height uint64) (resEndBlock types.ResponseEndBlock) {
	//log.Println("GlogChainApp.EndBlock", "height=", height, "size=", app.State.Size())
	return
}

