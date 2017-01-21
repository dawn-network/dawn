package main

import (
	"fmt"
	. "github.com/tendermint/go-common"
	"encoding/hex"
	"github.com/baabeetaa/glogchain/protocol"
	"log"
	"github.com/baabeetaa/glogchain/db"
	"github.com/tendermint/abci/types"
	"bytes"
	"encoding/gob"
	"github.com/tendermint/go-merkle"
	"encoding/json"
	dbm "github.com/tendermint/go-db"
	"github.com/tendermint/go-wire"
)

type GlogChainApp struct {
	db   		dbm.DB
	state 		merkle.Tree
	height 		uint64
	hashCount 	int
	txCount   	int
}

func NewGlogChainApp() *GlogChainApp {
	log.Println("NewGlogChainApp")

	//db := dbm.NewDB("state", "leveldb", ".")
	db := dbm.NewDB("state", dbm.GoLevelDBBackendStr, ".")

	lastBlock := LoadLastBlock(db)

	state := merkle.NewIAVLTree(0, db)
	state.Load(lastBlock.AppHash)

	log.Println("Loaded state", "block", lastBlock.Height, "root", state.Hash())

	return &GlogChainApp{ db: db, state: state}
}

func (app *GlogChainApp) Info() (resInfo types.ResponseInfo) {
	log.Println("GlogChainApp.Info")
	//return types.ResponseInfo{Data: Fmt("{\"hashes\":%v,\"txs\":%v}", app.hashCount, app.txCount)}

	resInfo = types.ResponseInfo{}
	resInfo.Data = Fmt("{\"hashes\":%v,\"txs\":%v}", app.hashCount, app.txCount)

	lastBlock := LoadLastBlock(app.db)
	resInfo.LastBlockHeight = lastBlock.Height
	resInfo.LastBlockAppHash = lastBlock.AppHash
	return resInfo
}

func (app *GlogChainApp) SetOption(key string, value string) (logstr string) {
	log.Println("GlogChainApp.SetOption")
	return ""
}

func (app *GlogChainApp) DeliverTx(tx []byte) types.Result {
	log.Println("GlogChainApp.DeliverTx")
	// tx is json string, need to convert to text and then parse into json object

	arr, err := hex.DecodeString(string(tx[:]))
	if (err != nil) {
		log.Println(err.Error())
		return types.ErrInternalError
	}

	jsonstring:= string(arr[:])


	// caculate hash of the tx
	operationEnvelope := protocol.OperationEnvelope{}
	err = json.Unmarshal([]byte(jsonstring), &operationEnvelope)
	if (err != nil) {
		log.Println(err.Error())
		return types.ErrEncodingError
	}
	opt_hash := protocol.Hash([]byte(operationEnvelope.Operation))


	app.state.Set(opt_hash, tx)



	///////////////////////////
	obj , err := protocol.UnMarshal(jsonstring)
	if (err != nil) {
		log.Println(err.Error())
		return types.ErrEncodingError
	}

	var buf bytes.Buffer        // Stand-in for a network connection
	enc := gob.NewEncoder(&buf) // Will write to network.
	dec := gob.NewDecoder(&buf) // Will read from network.

	err = enc.Encode(obj)
	if (err != nil) {
		log.Println(err.Error())
		return types.ErrInternalError
	}

	switch obj.(type) {  //v:=obj.(type) {
	case protocol.PostCreateOperation:
		//var tx protocol.PostCreateOperation
		//tx = v

		var post db.Post
		err = dec.Decode(&post)
		if err != nil {
			log.Println(err.Error())
			return types.ErrInternalError
		}

		err = db.CreatePost(post)
		if err != nil {
			log.Println(err.Error())
			return types.ErrInternalError
		}
	case protocol.PostEditOperation:
		//var tx protocol.PostEditOperation
		//tx = v

		var post db.Post
		err = dec.Decode(&post)
		if err != nil {
			log.Println(err.Error())
			return types.ErrInternalError
		}

		err = db.EditPost(post)
		if err != nil {
			log.Println(err.Error())
			return types.ErrInternalError
		}
	case protocol.AccountCreateOperation:
		//var tx protocol.AccountCreateOperation
		//tx = v

		var user db.User
		err = dec.Decode(&user)
		if err != nil {
			log.Println(err.Error())
			return types.ErrInternalError
		}

		err = db.CreateUser(user)
		if err != nil {
			log.Println(err.Error())
			return types.ErrInternalError
		}
	default:
	}

	app.txCount++
	return types.NewResult(types.CodeType_OK, tx, "DeliverTx OK")
}

func (app *GlogChainApp) CheckTx(tx []byte) types.Result {
	log.Println("GlogChainApp.CheckTx")
	dst := make([]byte, len(tx) * 2)
	hex.Encode(dst, tx)
	fmt.Println("CheckTx: ", dst)

	return types.OK
}

func (app *GlogChainApp) Commit() types.Result {
	log.Println("GlogChainApp.Commit")

	appHash := app.state.Save()

	log.Println("Saved state", "root", appHash)

	lastBlock := LastBlockInfo {
		Height:  app.height,
		AppHash: appHash, // this hash will be in the next block header
	}

	SaveLastBlock(app.db, lastBlock)

	return types.NewResultOK(appHash, "")

	//hash := app.state.Hash()
	//return types.NewResultOK(hash, "")
}

func (app *GlogChainApp) Query(query []byte) types.Result {
	log.Println("GlogChainApp.Query")
	return types.NewResultOK(nil, fmt.Sprintf("Query is not supported"))
}

func (app *GlogChainApp) InitChain(vals []*types.Validator) {
	log.Println("GlogChainApp.InitChain")
	//for _, plugin := range app.plugins.GetList() {
	//	plugin.InitChain(app.state, validators)
	//}
}

// TMSP::BeginBlock
func (app *GlogChainApp) BeginBlock(hash []byte, header *types.Header) {
	log.Println("GlogChainApp.BeginBlock")
	app.height = header.Height

	//for _, plugin := range app.plugins.GetList() {
	//	plugin.BeginBlock(app.state, height)
	//}
}

// TMSP::EndBlock
func (app *GlogChainApp) EndBlock(height uint64) (resEndBlock types.ResponseEndBlock) {
	log.Println("GlogChainApp.EndBlock", height)
	//app.height = height

	//for _, plugin := range app.plugins.GetList() {
	//	moreDiffs := plugin.EndBlock(app.state, height)
	//	diffs = append(diffs, moreDiffs...)
	//}
	return
}



//-----------------------------------------
// persist the last block info

var lastBlockKey = []byte("lastblock")

type LastBlockInfo struct {
	Height  uint64
	AppHash []byte
}

// Get the last block from the db
func LoadLastBlock(db dbm.DB) (lastBlock LastBlockInfo) {
	buf := db.Get(lastBlockKey)
	if len(buf) != 0 {
		r, n, err := bytes.NewReader(buf), new(int), new(error)
		wire.ReadBinaryPtr(&lastBlock, r, 0, n, err)
		if *err != nil {
			// DATA HAS BEEN CORRUPTED OR THE SPEC HAS CHANGED
			Exit(Fmt("Data has been corrupted or its spec has changed: %v\n", *err))
		}
		// TODO: ensure that buf is completely read.
	}

	return lastBlock
}

func SaveLastBlock(db dbm.DB, lastBlock LastBlockInfo) {
	log.Println("Saving block", "height", lastBlock.Height, "root", lastBlock.AppHash)
	buf, n, err := new(bytes.Buffer), new(int), new(error)
	wire.WriteBinary(lastBlock, buf, n, err)
	if *err != nil {
		// TODO
		PanicCrisis(*err)
	}
	db.Set(lastBlockKey, buf.Bytes())
}