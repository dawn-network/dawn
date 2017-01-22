package app

import (
	"fmt"
	. "github.com/tendermint/go-common"
	"encoding/hex"
	"log"
	"github.com/baabeetaa/glogchain/db"
	"github.com/tendermint/abci/types"
	"bytes"
	"encoding/gob"
	"github.com/tendermint/go-merkle"
	dbm "github.com/tendermint/go-db"
	"github.com/tendermint/go-wire"
)

type GlogChainApp struct {
	Db   		dbm.DB
	State 		merkle.Tree
	Height 		uint64
	TxCount   	uint64
}

func NewGlogChainApp() *GlogChainApp {
	log.Println("NewGlogChainApp")

	//db := dbm.NewDB("state", "leveldb", ".")
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

	lastBlock := LoadLastBlock(app.Db)
	resInfo.LastBlockHeight = lastBlock.Height
	resInfo.LastBlockAppHash = lastBlock.AppHash
	return resInfo
}

func (app *GlogChainApp) SetOption(key string, value string) (logstr string) {
	log.Println("GlogChainApp.SetOption", key, value)
	return Exec_SetOption(app, key, value)
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

	//// caculate hash of the tx
	//operationEnvelope := protocol.OperationEnvelope{}
	//err = json.Unmarshal([]byte(jsonstring), &operationEnvelope)
	//if (err != nil) {
	//	log.Println(err.Error())
	//	return types.ErrEncodingError
	//}
	//opt_hash := protocol.Hash([]byte(operationEnvelope.Operation))
	//app.state.Set(opt_hash, tx)

	///////////////////////////
	obj , err := UnMarshal(jsonstring)
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
	case PostCreateOperation:
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
	case PostEditOperation:
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
	case AccountCreateOperation:
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

		//////////////////////////
		// store to state
		var account Account
		account.PubKey, err = hex.DecodeString(user.Pubkey)
		if (err != nil) {
			log.Println(err.Error())
			return types.ErrInternalError
		}
		//copy(account.PubKey, barr)

		account.Sequence = 1
		account.Balance = 0

		err = TreeSaveAccount(app.State, account)
		if (err != nil) {
			log.Println(err.Error())
			return types.ErrInternalError
		}
	default:
	}

	app.TxCount++
	//return types.NewResult(types.CodeType_OK, tx, "DeliverTx OK")
	return types.OK
}

func (app *GlogChainApp) CheckTx(tx []byte) types.Result {
	log.Println("GlogChainApp.CheckTx")
	dst := make([]byte, len(tx) * 2)
	hex.Encode(dst, tx)
	fmt.Println("CheckTx: ", string(dst[:]))

	return types.OK
}

func (app *GlogChainApp) Commit() types.Result {
	log.Println("GlogChainApp.Commit")

	// try to not create block if no TX but doesnt work!!!
	//if (app.txCount <= 0) {
	//	return types.NewError(types.CodeType_InternalError, "No Tx in block")
	//}

	appHash := app.State.Save()
	log.Println("Saved state", "root", appHash)

	lastBlock := LastBlockInfo {
		Height:  app.Height,
		AppHash: appHash, // this hash will be in the next block header
	}

	SaveLastBlock(app.Db, lastBlock)
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
	log.Println("GlogChainApp.BeginBlock", hash)
	app.Height = header.Height
}

func (app *GlogChainApp) EndBlock(height uint64) (resEndBlock types.ResponseEndBlock) {
	log.Println("GlogChainApp.EndBlock", "height=", height, "size=", app.State.Size())
	return
}



//-----------------------------------------
// persist the last block info

var lastBlockKey = []byte("lastblock")

type LastBlockInfo struct {
	Height  uint64
	AppHash []byte
	TxCount uint64
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

////////////////////////
