package app

import (
	"encoding/json"
	"log"
	"fmt"
	"encoding/hex"
	"github.com/dawn-network/glogchain/gopressdb"
)

// In prototype, we'll use json because we don't need high performance and protocols will need to be change however.
// use dynamic json, more at http://eagain.net/articles/go-dynamic-json/
// should look at steem operations for references
// https://github.com/steemit/steem/blob/73a2b522e609925d444acfeb40264c5a0e682705/libraries/protocol/include/steemit/protocol/operations.hpp

type OperationEnvelope struct {
	Type 		string
	Operation 	string 		// json hex string of PostCreateOperation, PostEditOperation ...
	Signature 	string 		// crypto.SignatureEd25519 to the Operation, which is in json string
	Pubkey 		string 		// to verify and indentify who makes the transaction
	Fee		int64
}

////////////////////////////////////////
// Account

type AccountCreateOperation db.User

//type AccountUpdateOperation struct {
//	// need to define here
//}

////////////////////////////////////////
// Posting

type PostCreateOperation db.Post
type PostEditOperation db.Post

//type VoteOperation struct {
//	PostId 		string
//	Voter 		string
//}


////////////////////////////////////////
// crypto currency
type SendTokenOperation struct {
	ToAddress 	string
	Amount 		int64
}

func UnMarshal(jsonstring string) (env OperationEnvelope, returnObj interface{}, err error) {
	log.Println("UnMarshal", jsonstring)

	err = json.Unmarshal([]byte(jsonstring), &env)
	if (err != nil) {
		log.Println(err.Error())
		return
	}

	opt_arr, err := hex.DecodeString(env.Operation)
	if (err != nil) {
		log.Fatal(err)
		return
	}

	//opt_str := string(opt_arr)

	switch env.Type {
	case "AccountCreateOperation":
		var accountCreateOperation AccountCreateOperation
		if err = json.Unmarshal(opt_arr, &accountCreateOperation); err != nil {
			log.Fatal(err)
			return
		}
		returnObj = accountCreateOperation
		break
	case "SendTokenOperation":
		var sendTokenOperation SendTokenOperation
		if err = json.Unmarshal(opt_arr, &sendTokenOperation); err != nil {
			log.Fatal(err)
			return
		}
		returnObj = sendTokenOperation
		break
	case "PostCreateOperation":
		var postOperation PostCreateOperation
		if err = json.Unmarshal(opt_arr, &postOperation); err != nil {
			log.Fatal(err)
			return
		}
		returnObj = postOperation
		break
	case "PostEditOperation":
		var postOperation PostEditOperation
		if err = json.Unmarshal(opt_arr, &postOperation); err != nil {
			log.Fatal(err)
			return
		}
		returnObj = postOperation
		break
	case "VoteOperation":
		log.Fatalf("not support this type yet: %q", env.Type)
		err = fmt.Errorf("not support this type yet")
		return
	default:
		log.Fatalf("unknown Operation type: %q", env.Type)
		err = fmt.Errorf("unknown Operation type")
		return
	}
	return
}

//func Marshal() {
//	s := OperationEnvelope{
//		Type: "PostOperation",
//		Operation: PostOperation{
//			Title: "the Title",
//			Body: "the Body",
//			Author: "the Author",
//		},
//	}
//	buf, err := json.Marshal(s)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("%s\n", buf)
//
//	//c := OperationEnvelope{
//	//	Type: "cowbell",
//	//	Msg: Cowbell{
//	//		More: true,
//	//	},
//	//}
//	//buf, err = json.Marshal(c)
//	//if err != nil {
//	//	log.Fatal(err)
//	//}
//	//fmt.Printf("%s\n", buf)
//}

