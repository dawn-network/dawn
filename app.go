package main

import (
	"encoding/binary"
	"fmt"

	. "github.com/tendermint/go-common"
	"github.com/tendermint/tmsp/types"
)

type GlogChainApp struct {
	hashCount int
	txCount   int
}

func NewGlogChainApp() *GlogChainApp {
	return &GlogChainApp{}
}

func (app *GlogChainApp) Info() string {
	return Fmt("hashes:%v, txs:%v", app.hashCount, app.txCount)
}

func (app *GlogChainApp) SetOption(key string, value string) (log string) {
	return ""
}

func (app *GlogChainApp) AppendTx(tx []byte) types.Result {
	app.txCount += 1
	return types.OK
}

func (app *GlogChainApp) CheckTx(tx []byte) types.Result {
	return types.OK
}

func (app *GlogChainApp) Commit() types.Result {
	app.hashCount += 1

	if app.txCount == 0 {
		return types.OK
	} else {
		hash := make([]byte, 8)
		binary.BigEndian.PutUint64(hash, uint64(app.txCount))
		return types.NewResultOK(hash, "")
	}
}

func (app *GlogChainApp) Query(query []byte) types.Result {
	return types.NewResultOK(nil, fmt.Sprintf("Query is not supported"))
}
