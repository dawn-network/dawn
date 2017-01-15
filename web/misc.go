package web

import (
	"net/http"
	"log"
	"strings"
	"encoding/hex"
	"github.com/tendermint/go-crypto"
	"reflect"
	"errors"
)

// http://stackoverflow.com/questions/20170275/how-to-find-a-type-of-a-object-in-golang
func GetType(v interface{}) string {
	return reflect.TypeOf(v).String()
}

// http://stackoverflow.com/questions/18276173/calling-a-template-with-several-pipeline-parameters
func Dict(values ...interface{}) (map[string]interface{}, error) {
	if len(values)%2 != 0 {
		return nil, errors.New("invalid dict call")
	}

	dict := make(map[string]interface{}, len(values)/2)

	for i := 0; i < len(values); i+=2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, errors.New("dict keys must be strings")
		}
		dict[key] = values[i+1]

		// fix wrong type for value; expected int64; got int
		switch dict[key].(type) {
		case int:
			dict[key] = int64(dict[key].(int))
		}
	}

	return dict, nil
}

func StringCut(str string, n int) string {
	if (n < len(str)) {
		return str[:n]
	}

	return str
}

func GenerateKeyPair(w http.ResponseWriter, req *http.Request) {
	secret := req.FormValue("secret")

	// If method is GET serve an html login page
	if req.Method != "POST" {
		log.Println("PostCreateHandler GET")
		context := Context{Title: "account_genkeypair", Data: map[string]interface{}{ "secret": secret}}
		render(w, "account_genkeypair", context)
		return
	}

	//////////////////////////////////
	privKey := crypto.GenPrivKeyEd25519FromSecret([]byte(secret))

	log.Println("------------------------------------------------------------")
	log.Println("secret=\t\t\t" + secret)

	prikeystr := strings.ToUpper(hex.EncodeToString(privKey.Bytes()))
	prikeystr = prikeystr[2:len(prikeystr)]

	log.Println("_priKeyHex=\t\t\t" + prikeystr)

	//pubkeystr := strings.ToUpper(hex.EncodeToString(privKey.PubKey().Bytes()))
	//pubkeystr = pubkeystr[2:len(pubkeystr)]
	//log.Println("PubKey Hex=\t\t\t" + pubkeystr)

	log.Println("PubKey KeyString=\t" + privKey.PubKey().KeyString())
	log.Println("PubKey Address=\t\t" + strings.ToUpper(hex.EncodeToString(privKey.PubKey().Address())))

	render(w, "account_genkeypair",
		ActionResult{
			Status: "success",
			Message: "ok",
			Data: map[string]interface{}{
				"secret": secret,
				"PriKey": prikeystr,
				"PubKey": privKey.PubKey().KeyString(),
				"Address": strings.ToUpper(hex.EncodeToString(privKey.PubKey().Address())),
			},
		})
}