package main

import (
	"fmt"
	. "github.com/tendermint/go-common"
	"encoding/hex"
	"github.com/baabeetaa/glogchain/protocol"
	"log"
	"github.com/baabeetaa/glogchain/db"
	"github.com/tendermint/abci/types"
)

type GlogChainApp struct {
	hashCount int
	txCount   int
}

func NewGlogChainApp() *GlogChainApp {
	return &GlogChainApp{}
}

func (app *GlogChainApp) Info() types.ResponseInfo {
	return types.ResponseInfo{Data: Fmt("{\"hashes\":%v,\"txs\":%v}", app.hashCount, app.txCount)}
}

func (app *GlogChainApp) SetOption(key string, value string) (log string) {
	return ""
}

func (app *GlogChainApp) DeliverTx(tx []byte) types.Result {
	// tx is json string, need to convert to text and then parse into json object
	jsonstring := string(tx[:])

	obj , err := protocol.UnMarshal(jsonstring)

	if err != nil {
		log.Fatal(err)
		return types.ErrEncodingError
	}

	switch v:=obj.(type) {
	case protocol.PostCreateOperation:
		var tx protocol.PostCreateOperation
		tx = v

		err = db.CreatePost(tx.Post)
		if err != nil {
			log.Println("AppendTx CreatePost", err)
			return types.ErrInternalError
		}
	case protocol.PostEditOperation:
		var tx protocol.PostEditOperation
		tx = v

		err = db.EditPost(tx.Post)
		if err != nil {
			log.Println("AppendTx EditPost", err)
			return types.ErrInternalError
		}
	case protocol.AccountCreateOperation:
		var tx protocol.AccountCreateOperation
		tx = v

		err = db.CreateUser(tx.User)
		if err != nil {
			log.Println("AppendTx CreateUser", err)
			return types.ErrInternalError
		}
	default:
	}

	//return types.OK
	return types.NewResult(types.CodeType_OK, tx, "AppendTx OK")
}

func (app *GlogChainApp) CheckTx(tx []byte) types.Result {
	dst := make([]byte, len(tx) * 2)
	hex.Encode(dst, tx)
	fmt.Println("CheckTx: ", dst)

	return types.OK
}

func (app *GlogChainApp) Commit() types.Result {
	app.hashCount += 1

	//if app.txCount == 0 {
		return types.OK
	//} else {
	//	hash := make([]byte, 8)
	//	binary.BigEndian.PutUint64(hash, uint64(app.txCount))
	//	return types.NewResultOK(hash, "")
	//}
}

func (app *GlogChainApp) Query(query []byte) types.Result {
	return types.NewResultOK(nil, fmt.Sprintf("Query is not supported"))
}
