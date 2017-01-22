package app

import (
	"log"
	"github.com/tendermint/go-merkle"
)

///////////////////////////////////////////
//-----------------------------------------
// account to store on merkle tree (key=account address, value=account)
type Account struct {
	PubKey   	[]byte //[32]byte , crypto.PubKeyEd25519
	Sequence 	int64
	Balance  	int64
}



func TreeSaveAccount(state merkle.Tree, acc Account) error  {
	pubkey, err := GetPubKeyFromBytes(acc.PubKey)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	raw, err := ToBytes(acc)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	state.Set(pubkey.Address(), raw)

	return nil
}
