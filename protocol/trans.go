package protocol

import (
	"encoding/json"
	"log"
	"fmt"
)


// In prototype, we'll use json because we don't need high performance and protocols will need to be change however.
// use dynamic json, more at http://eagain.net/articles/go-dynamic-json/
// should look at steem operations for references
// https://github.com/steemit/steem/blob/73a2b522e609925d444acfeb40264c5a0e682705/libraries/protocol/include/steemit/protocol/operations.hpp

type OperationEnvelope struct {
	Type 		string
	Operation 	interface{}
}

type AccountCreateOperation struct {
	Fee		float64
	PubKey 		string
}

type AccountUpdateOperation struct {
	// need to define here
}

type PostOperation struct {
	Fee		float64
	Title 		string
	Body 		string
	Author 		string
}

type VoteOperation struct {
	Fee		float64
	PostId 		string
	Voter 		string
}


func UnMarshal(jsonstring string) (interface{}, error) {
	var returnObj interface{}

	var operation json.RawMessage
	env := OperationEnvelope{
		Operation: &operation,
	}

	if err := json.Unmarshal([]byte(jsonstring), &env); err != nil {
		log.Fatal(err)
		return nil, err
	}

	switch env.Type {
	case "AccountCreateOperation":
		var accountCreateOperation AccountCreateOperation
		if err := json.Unmarshal(operation, &accountCreateOperation); err != nil {
			log.Fatal(err)
			return nil, err
		}
		returnObj = accountCreateOperation
	case "PostOperation":
		var postOperation PostOperation
		if err := json.Unmarshal(operation, &postOperation); err != nil {
			log.Fatal(err)
			return nil, err
		}
		//var desc string = s.Description
		//fmt.Println(desc)

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

func Marshal() {
	s := OperationEnvelope{
		Type: "PostOperation",
		Operation: PostOperation{
			Title: "the Title",
			Body: "the Body",
			Author: "the Author",
		},
	}
	buf, err := json.Marshal(s)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", buf)

	//c := OperationEnvelope{
	//	Type: "cowbell",
	//	Msg: Cowbell{
	//		More: true,
	//	},
	//}
	//buf, err = json.Marshal(c)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Printf("%s\n", buf)
}