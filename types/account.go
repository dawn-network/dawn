package types

///////////////////////////////////////////
//-----------------------------------------
// account to store on merkle tree (key=account address, value=account)
type Account struct {
	PubKey   	[]byte 	//[32]byte , crypto.PubKeyEd25519
	Sequence 	int64
	Balance  	int64
}

