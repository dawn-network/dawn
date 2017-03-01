package app

import (
	tm_types "github.com/tendermint/abci/types"
	"log"
	"encoding/hex"
	"bytes"
	"encoding/gob"
	//"github.com/dawn-network/glogchain/core/db"
	"github.com/dawn-network/glogchain/types"
)

func Exec_DeliverTx(tx []byte) tm_types.Result {
	// tx is json string, need to convert to text and then parse into json object

	var err error

	arr, err := hex.DecodeString(string(tx[:]))
	if (err != nil) {
		log.Println(err.Error())
		return tm_types.ErrInternalError
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
		return tm_types.ErrEncodingError
	}


	bPubKey, err := hex.DecodeString(env.Pubkey)
	if (err != nil) {
		log.Println(err.Error())
		return tm_types.ErrEncodingError
	}

	PubKey, err := GetPubKeyFromBytes(bPubKey)
	if (err != nil) {
		log.Println(err.Error())
		return tm_types.ErrInternalError
	}

	Address := PubKey.Address() // address of the user who make the transaction

	var buf bytes.Buffer        // Stand-in for a network connection
	enc := gob.NewEncoder(&buf) // Will write to network.
	dec := gob.NewDecoder(&buf) // Will read from network.

	err = enc.Encode(obj)
	if (err != nil) {
		log.Println(err.Error())
		return tm_types.ErrInternalError
	}

	switch obj.(type) {
	case AccountCreateOperation:
		var user types.User
		err = dec.Decode(&user)
		if err != nil {
			log.Println(err.Error())
			return tm_types.ErrInternalError
		}

		//err = db.CreateUser(user)
		//if err != nil {
		//	log.Println(err.Error())
		//	return tm_types.ErrInternalError
		//}

		//////////////////////////
		// store to state
		var account types.Account
		account.PubKey, err = hex.DecodeString(user.Pubkey)
		if (err != nil) {
			log.Println(err.Error())
			return tm_types.ErrInternalError
		}
		//copy(account.PubKey, barr)

		account.Sequence = 1
		account.Balance = 0

		err = TreeSaveAccount(GlogGlobal.GlogApp.State, account)
		if (err != nil) {
			log.Println(err.Error())
			return tm_types.ErrInternalError
		}
		break
	case SendTokenOperation:
		var sendtoken SendToken
		sendtoken.From = Address

		opt, ok := obj.(SendTokenOperation)
		if (!ok) {
			log.Println(err.Error())
			return tm_types.ErrInternalError
		}

		ToAddress, err := hex.DecodeString(opt.ToAddress)
		if (err != nil) {
			log.Println(err.Error())
			return tm_types.ErrEncodingError
		}

		sendtoken.To = ToAddress
		sendtoken.Amount = opt.Amount

		err = Exec_SendToken(GlogGlobal.GlogApp.State, sendtoken)
		if (err != nil) {
			log.Println(err.Error())
			return tm_types.ErrEncodingError
		}

		break
	case PostCreateOperation:
		var post types.Post
		err = dec.Decode(&post)
		if err != nil {
			log.Println(err.Error())
			return tm_types.ErrInternalError
		}

		//err = db.CreatePost(post)
		//if err != nil {
		//	log.Println(err.Error())
		//	return tm_types.ErrInternalError
		//}
		break
	case PostEditOperation:
		var post types.Post
		err = dec.Decode(&post)
		if err != nil {
			log.Println(err.Error())
			return tm_types.ErrInternalError
		}

		//err = db.EditPost(post)
		//if err != nil {
		//	log.Println(err.Error())
		//	return tm_types.ErrInternalError
		//}
		break
	case CommentCreateOperation:
		var comment types.Comment
		err = dec.Decode(&comment)
		if err != nil {
			log.Println(err.Error())
			return tm_types.ErrInternalError
		}

		//err = db.CreateComment(comment)
		//if err != nil {
		//	log.Println(err.Error())
		//	return tm_types.ErrInternalError
		//}
		break
	default:
		log.Println("Unknow Operation")
	}

	GlogGlobal.GlogApp.TxCount++
	//return types.NewResult(types.CodeType_OK, tx, "DeliverTx OK")
	return tm_types.OK
}