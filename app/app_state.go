package app

import (
	"bytes"
	"log"
	dbm "github.com/tendermint/go-db"
	"github.com/tendermint/go-wire"
	. "github.com/tendermint/go-common"
)

//-----------------------------------------
// persist the last block info

var lastBlockKey = []byte("lastblock")

type LastBlockInfo struct {
	Height  	uint64
	AppHash 	[]byte
	TxCount 	uint64
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
