package service

import (
	"math/rand"
	"log"
	"io/ioutil"
	"encoding/hex"
	//"github.com/dawn-network/glogchain/app"
	"net/http"
	"github.com/dawn-network/glogchain/types"
)

//func categories_normalize(jsonstr string) (string, error) {
//	cats_string := []string{}
//	json.Unmarshal([]byte(jsonstr), &cats_string)
//
//	//items := []Category{}
//	for i, item := range cats_string {
//	//	cat := Category{ item, 0 }
//	//	items = append(items, cat)
//
//		item = strings.ToLower(item)
//		item = strings.TrimSpace(item)
//		cats_string[i] = item
//	}
//
//	return items, nil
//
//
//}


var letters = []rune("abcdefghijklmnopqrstuvwxyz0123456789")

/**
 * Generate random string
 * http://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-golang
 */
func RandSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}


/////////////////

func TM_broadcast_tx_commit(tx []byte) {
	log.Println("tx_json=", string(tx[:]))

	tx_json_hex := make([]byte, len(tx) * 2)
	hex.Encode(tx_json_hex, tx)
	log.Println("tx_json_hex", string(tx_json_hex[:]))
	//
	data := string(tx_json_hex[:])


	var url_request string = types.GlogchainConfigGlobal.TmRpcLaddr + "/broadcast_tx_commit?tx=%22" + data + "%22"
	log.Println("TM_broadcast_tx_commit url_request: %#v\n", url_request)
	resp, err := http.Get(url_request)
	if err != nil {
		log.Println("TM_broadcast_tx_commit http.Get error", err.Error())
		return;
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("TM_broadcast_tx_commit ioutil.ReadAll error", err.Error())
		return;
	}
	json_response_string := string(body[:])
	log.Println("TM_broadcast_tx_commit json_response_string:", json_response_string)
}