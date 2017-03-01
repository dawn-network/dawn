package web

import (
	"log"
	"encoding/hex"
	"net/http"
	"github.com/tendermint/go-crypto"
	"github.com/dawn-network/glogchain/core/app"
	"strings"
	"strconv"
	"encoding/json"
	"github.com/dawn-network/glogchain/core/service"
)

func WalletViewHandler(w http.ResponseWriter, req *http.Request) {
	session := GetSession(req)

	private_key, ok := session.Values["private_key"].(crypto.PrivKeyEd25519)
	if (!ok) {
		log.Println("Can not get session private_key")
		return
	}

	acc, err := app.TreeGetAccount(app.GlogGlobal.GlogApp.State, private_key.PubKey().Address())
	if (err != nil) {
		log.Println("Can not get account from app state")
		return
	}

	render(w, "wallet_view", ActionResult{Status: "ok", Message: "msg ok",
		Data: map[string]interface{}{
			"Address": strings.ToUpper(hex.EncodeToString(private_key.PubKey().Address())),
			"PubKey": strings.ToUpper(hex.EncodeToString(acc.PubKey)),
			"Sequence": acc.Sequence,
			"Balance": acc.Balance,
		}})
}

func WalletSendTokenHandler(w http.ResponseWriter, req *http.Request) {
	session := GetSession(req)

	pToAddress := req.FormValue("ToAddress")
	pAmount := req.FormValue("Amount")

	// If method is GET serve an html
	if req.Method != "POST" {
		log.Println("wallet_sendtoken GET")
		context := Context{Title: "wallet_sendtoken",
			Data: map[string]interface{}{
				"ToAddress": pToAddress,
				"Amount": pAmount,
			}}
		render(w, "wallet_sendtoken", context)
		return
	}

	///////////////////////////////
	// form validate
	if (len(pToAddress) != 40) {
		render(w, "wallet_sendtoken",
			ActionResult{
				Status: "error",
				Message: "Invalid ToAddress",
				Data: map[string]interface{}{
					"ToAddress": pToAddress,
					"Amount": pAmount,
				},
			})
		return
	}

	ToAddress, err := hex.DecodeString(pToAddress)
	if (err != nil) {
		render(w, "wallet_sendtoken",
			ActionResult{
				Status: "error",
				Message: "Invalid ToAddress",
				Data: map[string]interface{}{
					"ToAddress": pToAddress,
					"Amount": pAmount,
				},
			})
		return
	}


	Amount, err := strconv.ParseInt(pAmount, 10, 64)
	if (err != nil) {
		render(w, "wallet_sendtoken",
			ActionResult{
				Status: "error",
				Message: "Amount is not a number",
				Data: map[string]interface{}{
					"ToAddress": pToAddress,
					"Amount": pAmount,
				},
			})
		return
	}

	log.Println("ToAddress=", hex.EncodeToString(ToAddress), "Amount=", Amount)


	opt := app.SendTokenOperation {
		ToAddress: strings.ToUpper(hex.EncodeToString(ToAddress)),
		Amount: Amount,
	}

	opt_arr, err := json.Marshal(opt)
	if (err != nil) {
		render(w, "wallet_sendtoken",
			ActionResult{
				Status: "error",
				Message: err.Error(),
				Data: map[string]interface{}{
					"ToAddress": pToAddress,
					"Amount": pAmount,
				},
			})
		return
	}
	opt_str := strings.ToUpper(hex.EncodeToString(opt_arr))

	// sign the transaction
	private_key, ok := session.Values["private_key"].(crypto.PrivKeyEd25519)
	if (!ok) {
		log.Fatal(err.Error())
		return
	}
	sign := private_key.Sign([]byte(opt_str))
	sign_str := strings.ToUpper(hex.EncodeToString(sign.Bytes()))
	sign_str = sign_str[2:len(sign_str)]

	tx := app.OperationEnvelope {
		Type: "SendTokenOperation",
		Operation: opt_str,
		Signature: sign_str,
		Pubkey: private_key.PubKey().KeyString(),
		Fee: 0,
	}

	byte_arr, err := json.Marshal(tx)
	if err != nil {
		log.Fatal("PostEditHandler Cannot encode to JSON ", err)
		return
	}

	//tx_json := string(byte_arr[:])
	//log.Println("PostEditHandler tx_json=", tx_json)
	//
	//tx_json_hex := make([]byte, len(tx_json) * 2)
	//hex.Encode(tx_json_hex, []byte(tx_json))
	//log.Println("PostEditHandler tx_json_hex", string(tx_json_hex[:]))
	//
	//service.TM_broadcast_tx_commit(string(tx_json_hex[:]))
	service.TM_broadcast_tx_commit(byte_arr)

	render(w, "wallet_sendtoken",
		ActionResult{
			Status: "success",
			Message: "send OK!",
			Data: map[string]interface{}{
				"ToAddress": pToAddress,
				"Amount": pAmount,
			},
		})

	//receiver, err := app.TreeGetAccount(app.GlogGlobal.GlogApp.State, tx.To)
	//if (err != nil) {
	//	return err
	//}

}