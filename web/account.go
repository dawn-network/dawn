package web

import (
	"strings"
	"encoding/hex"
	"bytes"
	"encoding/binary"
	"log"
	"github.com/baabeetaa/glogchain/protocol"
	"encoding/json"
	"github.com/baabeetaa/glogchain/service"
	"net/http"
	"github.com/baabeetaa/glogchain/config"
	"github.com/tendermint/go-crypto"
)

func AccountCreate(w http.ResponseWriter, req *http.Request) {
	// If method is GET serve an html
	if req.Method != "POST" {
		context := Context{Title: "Welcome!"}
		context.Static = config.GlogchainConfigGlobal.WebRootDir + "/static/"
		context.Data = map[string]interface{}{ "username": "", "prikey": ""}
		render(w, "account_create", context)
		return
	}

	username := req.FormValue("username")
	prikey := req.FormValue("prikey")

	log.Println("AccountCreateHandler", "username", username)
	log.Println("AccountCreateHandler", "prikey", prikey)

	if (len(username) < 6) {
		render(w, "account_create", ActionResult{Status: "error", Message: "username must be at least 6 characters", Data: map[string]interface{}{ "username": username, "prikey": prikey}})
		return
	}

	if (len(prikey) != 128) {
		render(w, "account_create", ActionResult{Status: "error", Message: "PubKey must be 32 bytes in Hex String ( 64 characters)", Data: map[string]interface{}{ "username": username, "prikey": prikey}})
		return
	}


	prikey = strings.ToUpper(prikey)
	byte_arr, err := hex.DecodeString(prikey)
	if (err != nil) {
		render(w, "account_create", ActionResult{Status: "error", Message: err.Error(), Data: map[string]interface{}{ "username": username, "prikey": prikey}})
		return
	}

	buf := &bytes.Buffer{}
	err = binary.Write(buf, binary.BigEndian, byte_arr)
	if err != nil {
		render(w, "account_create", ActionResult{Status: "error", Message: err.Error(), Data: map[string]interface{}{ "username": username, "prikey": prikey}})
		return
	}

	var key crypto.PrivKeyEd25519
	binary.Read(buf, binary.BigEndian, &key);

	log.Println("AccountCreateHandler", "PubKey", key.PubKey().KeyString())
	var address string = strings.ToUpper(hex.EncodeToString(key.PubKey().Address()))
	log.Println("AccountCreateHandler Address=\t\t" + address)


	//////////////////////////////////////
	// build the transaction
	var opt protocol.AccountCreateOperation
	opt.ID = address
	opt.Username = username
	opt.Pubkey = key.PubKey().KeyString()
	opt.UserRegistered = "2017-01-06 09:00:28"
	opt.DisplayName = username

	opt_arr, err := json.Marshal(opt)
	if (err != nil) {
		render(w, "account_create", ActionResult{Status: "error", Message: err.Error(), Data: map[string]interface{}{ "username": username, "prikey": prikey}})
		return
	}
	opt_str := strings.ToUpper(hex.EncodeToString(opt_arr))

	// sign the transaction
	sign := key.Sign(opt_arr)
	sign_str := strings.ToUpper(hex.EncodeToString(sign.Bytes()))
	sign_str = sign_str[2:len(sign_str)]

	tx := protocol.OperationEnvelope{
		Type: "AccountCreateOperation",
		Operation: opt_str,
		Signature: sign_str,
		Pubkey: key.PubKey().KeyString(),
		Fee: 0,
	}

	byte_arr, err = json.Marshal(tx)
	if err != nil {
		log.Fatal("AccountCreateHandler Cannot encode to JSON ", err)
	}

	tx_json := string(byte_arr[:])
	log.Println("AccountCreateHandler tx_json=", tx_json)

	tx_json_hex := make([]byte, len(tx_json) * 2)
	hex.Encode(tx_json_hex, []byte(tx_json))
	log.Println("AccountCreateHandler tx_json_hex", string(tx_json_hex[:]))

	service.TM_broadcast_tx_commit(string(tx_json_hex[:]))
	render(w, "account_create", ActionResult{Status: "success", Message: "ok", Data: map[string]interface{}{ "username": username, "prikey": prikey}})
}



