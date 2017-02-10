package web

import (
	"strings"
	"encoding/hex"
	"log"
	"encoding/json"
	"github.com/dawn-network/glogchain/service"
	"net/http"
	"github.com/dawn-network/glogchain/app"
)

func AccountCreate(w http.ResponseWriter, req *http.Request) {
	// If method is GET serve an html
	if req.Method != "POST" {
		context := Context{Title: "Welcome!"}
		context.Static = app.GlogchainConfigGlobal.WebRootDir + "/static/"
		context.Data = map[string]interface{}{ "username": "", "prikey": ""}
		render(w, "account_create", context)
		return
	}

	username := req.FormValue("username")
	prikey_str := req.FormValue("prikey")

	log.Println("username", username)
	log.Println("prikey", prikey_str)

	if (len(username) < 3) {
		render(w, "account_create", ActionResult{Status: "error", Message: "username must be at least 3 characters", Data: map[string]interface{}{ "username": username, "prikey": prikey_str}})
		return
	}

	if (len(prikey_str) != 128) {
		render(w, "account_create", ActionResult{Status: "error", Message: "PubKey must be 32 bytes in Hex String ( 64 characters)", Data: map[string]interface{}{ "username": username, "prikey": prikey_str}})
		return
	}


	prikey_str = strings.ToUpper(prikey_str)
	prikey_byte_arr, err := hex.DecodeString(prikey_str)
	if (err != nil) {
		render(w, "account_create", ActionResult{Status: "error", Message: err.Error(), Data: map[string]interface{}{ "username": username, "prikey": prikey_str}})
		return
	}

	private_key, err := app.GetPrivateKeyFromBytes(prikey_byte_arr)
	if err != nil {
		render(w, "account_create", ActionResult{Status: "error", Message: err.Error(), Data: map[string]interface{}{ "username": username, "prikey": prikey_str}})
		return
	}

	log.Println("PubKey", private_key.PubKey().KeyString())
	var address string = strings.ToUpper(hex.EncodeToString(private_key.PubKey().Address()))
	log.Println("Address=\t\t" + address)


	//////////////////////////////////////
	// build the transaction
	var opt app.AccountCreateOperation
	opt.ID = address
	opt.Username = username
	opt.Pubkey = private_key.PubKey().KeyString()
	opt.UserRegistered = "2017-01-06 09:00:28"
	opt.DisplayName = username

	opt_arr, err := json.Marshal(opt)
	if (err != nil) {
		render(w, "account_create", ActionResult{Status: "error", Message: err.Error(), Data: map[string]interface{}{ "username": username, "prikey": prikey_str}})
		return
	}
	opt_str := strings.ToUpper(hex.EncodeToString(opt_arr))

	// sign the transaction
	sign := private_key.Sign([]byte(opt_str))
	sign_str := strings.ToUpper(hex.EncodeToString(sign.Bytes()))
	sign_str = sign_str[2:len(sign_str)]

	// test verifying
	log.Println("vefify", private_key.PubKey().VerifyBytes(opt_arr, sign))

	tx := app.OperationEnvelope{
		Type: "AccountCreateOperation",
		Operation: opt_str,
		Signature: sign_str,
		Pubkey: private_key.PubKey().KeyString(),
		Fee: 0,
	}

	byte_arr, err := json.Marshal(tx)
	if err != nil {
		log.Fatal("Cannot encode to JSON ", err)
		return
	}

	log.Println("tx_json=", string(byte_arr[:]))

	tx_json_hex := make([]byte, len(byte_arr) * 2)
	hex.Encode(tx_json_hex, byte_arr)
	log.Println("tx_json_hex", string(tx_json_hex[:]))

	service.TM_broadcast_tx_commit(string(tx_json_hex[:]))
	render(w, "account_create", ActionResult{Status: "success", Message: "ok", Data: map[string]interface{}{ "username": username, "prikey": prikey_str}})
}

