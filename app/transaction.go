package app

import (
	"github.com/tendermint/go-merkle"
	"log"
	"errors"
	"bytes"
	"github.com/tendermint/go-wire"
	"github.com/dawn-network/glogchain/types"
	"github.com/dawn-network/glogchain/service"
)




type SendToken struct {
	From		[]byte 	// sender address
	To	 	[]byte 	// receiver address
	Amount 		int64 	// how much
}

func TreeSaveAccount(state merkle.Tree, acc types.Account) error  {
	pubkey, err := service.GetPubKeyFromBytes(acc.PubKey)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	raw, err := service.StructToBytes(acc)
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