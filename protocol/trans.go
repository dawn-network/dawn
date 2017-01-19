package protocol

import (
	"github.com/baabeetaa/glogchain/db"
	"encoding/json"
	"log"
	"fmt"
	"encoding/hex"
)


// In prototype, we'll use json because we don't need high performance and protocols will need to be change however.
// use dynamic json, more at http://eagain.net/articles/go-dynamic-json/
// should look at steem operations for references
// https://github.com/steemit/steem/blob/73a2b522e609925d444acfeb40264c5a0e682705/libraries/protocol/include/steemit/protocol/operations.hpp

type OperationEnvelope struct {
	Type 		string
	Operation 	string 		// json hex string of PostCreateOperation, PostEditOperation ...
	Signature 	string 		// crypto.SignatureEd25519 to the Operation, which is in json
	Pubkey 		string 		// to verify and indentify who makes the transaction
	Fee		float64
}

type AccountCreateOperation db.User

//type AccountUpdateOperation struct {
//	// need to define here
//}

type PostCreateOperation db.Post

type PostEditOperation db.Post

//type VoteOperation struct {
//	PostId 		string
//	Voter 		string
//}


func UnMarshal(jsonstring string) (interface{}, error) {
	var returnObj interface{}

	//var operation json.RawMessage
	env := OperationEnvelope{
		//Operation: &operation,
	}

	if err := json.Unmarshal([]byte(jsonstring), &env); err != nil {
		log.Fatal(err)
		return nil, err
	}


	opt_arr, err := hex.DecodeString(env.Operation)
	if (err != nil) {
		log.Fatal(err)
		return nil, err
	}

	//opt_str := string(opt_arr)

	switch env.Type {
	case "AccountCreateOperation":
		var accountCreateOperation AccountCreateOperation
		if err := json.Unmarshal(opt_arr, &accountCreateOperation); err != nil {
			log.Fatal(err)
			return nil, err
		}
		returnObj = accountCreateOperation
	case "PostCreateOperation":
		var postOperation PostCreateOperation
		if err := json.Unmarshal(opt_arr, &postOperation); err != nil {
			log.Fatal(err)
			return nil, err
		}

		returnObj = postOperation
	case "PostEditOperation":
		var postOperation PostEditOperation
		if err := json.Unmarshal(opt_arr, &postOperation); err != nil {
			log.Fatal(err)
			return nil, err
		}

		returnObj = postOperation
	case "VoteOperation":
		log.Fatalf("not support this type yet: %q", env.Type)
		return nil, fmt.Errorf("not support this type yet")
	default:
		log.Fatalf("unknown Operation type: %q", env.Type)
		return nil, fmt.Errorf("unknown Operation type")
	}


	return returnObj, nil
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