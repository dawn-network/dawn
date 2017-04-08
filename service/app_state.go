package service

import (
	"bytes"
	"log"
	dbm "github.com/tendermint/go-db"
	"github.com/tendermint/go-wire"
	. "github.com/tendermint/go-common"
	"encoding/hex"
	"github.com/tendermint/go-merkle"
	"github.com/dawn-network/glogchain/types"
	"errors"
)

//-----------------------------------------
// persist the last block info

const lastBlockKey = "lastblock"

type LastBlockInfo struct {
	Height  	uint64
	AppHash 	[]byte
	TxCount 	uint64
}

// Get the last block from the db
func LoadLastBlock(db dbm.DB) (lastBlock LastBlockInfo) {
	buf := db.Get([]byte(lastBlockKey))
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
	log.Println("Saving block", "height", lastBlock.Height, "root", hex.EncodeToString(lastBlock.AppHash))
	buf, n, err := new(bytes.Buffer), new(int), new(error)
	wire.WriteBinary(lastBlock, buf, n, err)
	if *err != nil {
		// TODO
		PanicCrisis(*err)
	}
	db.Set([]byte(lastBlockKey), buf.Bytes())
}


//////////////////////

type SendToken struct {
	From		[]byte 	// sender address
	To	 	[]byte 	// receiver address
	Amount 		int64 	// how much
}

func TreeSaveAccount(state merkle.Tree, acc types.Account) error  {
	pubkey, err := GetPubKeyFromBytes(acc.PubKey)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	raw, err := StructToBytes(acc)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	state.Set(pubkey.Address(), raw)

	return nil
}

func TreeGetAccount(state merkle.Tree, key []byte) (acc types.Account, err error) {
	//var index int
	var value []byte
	var exists bool

	// get account
	_, value, exists = state.Get(key)
	if (!exists) {
		err = errors.New("no sender found")
		return
	}

	//var ok bool
	//acc, ok = BytesToStruct(acc, value).(Account)
	//if (!ok) {
	//	err = errors.New("Cannot read Account")
	//	return
	//}

	r, n := bytes.NewReader(value), new(int)
	wire.ReadBinaryPtr(&acc, r, 0, n, &err)
	if (err != nil) {
		return
	}

	return acc, err
}


func Exec_SendToken(state merkle.Tree, tx SendToken) error {
	sender, err := TreeGetAccount(state, tx.From)
	if (err != nil) {
		return err
	}

	receiver, err := TreeGetAccount(state, tx.To)
	if (err != nil) {
		return err
	}

	if (tx.Amount <= 0) {
		return errors.New("Amout value is invalid")
	}

	if (sender.Balance < tx.Amount) {
		return errors.New("Amout value is invalid")
	}

	// TODO: more validation
	sender.Balance -= tx.Amount
	receiver.Balance += tx.Amount

	sender.Sequence++
	receiver.Sequence++

	TreeSaveAccount(state, sender)
	TreeSaveAccount(state, receiver)

	return nil // everything OK, return nil
}