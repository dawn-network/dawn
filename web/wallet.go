package web

import (
	"log"
	"encoding/hex"
	"net/http"
	"github.com/tendermint/go-crypto"
	"github.com/baabeetaa/glogchain/app"
	"strings"
	"strconv"
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
	if err != nil {
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


	//receiver, err := app.TreeGetAccount(app.GlogGlobal.GlogApp.State, tx.To)
	//if (err != nil) {
	//	return err
	//}

}