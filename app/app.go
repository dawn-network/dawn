package app

import (
	"fmt"
	"encoding/hex"
	"log"
	"github.com/baabeetaa/glogchain/db"
	"github.com/tendermint/abci/types"
	"bytes"
	"encoding/gob"
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
	return Exec_SetOption(app, key, value)
}

func (app *GlogChainApp) DeliverTx(tx []byte) types.Result {
	log.Println("GlogChainApp.DeliverTx")
	// tx is json string, need to convert to text and then parse into json object

	var err error

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
	env, obj , err := UnMarshal(jsonstring)
	if (err != nil) {
		log.Println(err.Error())
		return types.ErrEncodingError
	}


	bPubKey, err := hex.DecodeString(env.Pubkey)
	if (err != nil) {
		log.Println(err.Error())
		return types.ErrEncodingError
	}

	PubKey, err := GetPubKeyFromBytes(bPubKey)
	if (err != nil) {
		log.Println(err.Error())
		return types.ErrInternalError
	}

	Address := PubKey.Address() // address of the user who make the transaction

	var buf bytes.Buffer        // Stand-in for a network connection
	enc := gob.NewEncoder(&buf) // Will write to network.
	dec := gob.NewDecoder(&buf) // Will read from network.

	err = enc.Encode(obj)
	if (err != nil) {
		log.Println(err.Error())
		return types.ErrInternalError
	}

	switch obj.(type) {  //v:=obj.(type) {
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
		break
	case SendTokenOperation:
		var sendtoken SendToken
		sendtoken.From = Address

		opt, ok := obj.(SendTokenOperation)
		if (!ok) {
			log.Println(err.Error())
			return types.ErrInternalError
		}

		ToAddress, err := hex.DecodeString(opt.ToAddress)
		if (err != nil) {
			log.Println(err.Error())
			return types.ErrEncodingError
		}

		sendtoken.To = ToAddress
		sendtoken.Amount = opt.Amount

		err = Exec_SendToken(GlogGlobal.GlogApp.State, sendtoken)
		if (err != nil) {
			log.Println(err.Error())
			return types.ErrEncodingError
		}

		break
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
		break
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
		break
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

