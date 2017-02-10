package web

import (
	"net/http"
	"github.com/dawn-network/glogchain/db"
	"time"
	"github.com/dawn-network/glogchain/app"
	"golang.org/x/crypto/ripemd160"
	"fmt"
	"github.com/dawn-network/glogchain/service"
	"strings"
	"encoding/hex"
	"encoding/json"
	"log"
	"github.com/tendermint/go-crypto"
)

func CommentCreateHandler(w http.ResponseWriter, req *http.Request) {
	var opt app.CommentCreateOperation
	var post db.Post
	var cm_parent db.Comment
	var err error

	session := GetSession(req)
	user, ok := session.Values["user"].(db.User)
	if (!ok) {
		panic("wtf?? session user is nil")
		return
	}

	p := req.FormValue("p") 	// postid
	c := req.FormValue("c") 	// parent comment id
	CommentContent := req.FormValue("CommentContent")
	opt.Content = CommentContent

	post, err = db.GetPost(p)
	opt.PostID = post.ID
	if (err != nil) {
		render(w, "comment_create", ActionResult{Status: "error", Message: "No Post found", Data: opt })
		return
	}

	if (c != "") {
		cm_parent, err = db.GetComment(c)
		opt.Parent = cm_parent.ID
		if (err != nil) {
			render(w, "comment_create", ActionResult{ Status: "error", Message: "No Comment found", Data: opt})
			return
		}
	}

	///////////////////////////////////////////
	if (req.Method != "POST") {
		render(w, "comment_create", ActionResult{ Status: "success", Message: "ok", Data: opt})
		return
	}

	///////////////////////////////////////////
	// Handler POST here
	timeNow := time.Now()

	opt.Author = user.ID 				// take from session
	opt.Date = timeNow.String() 			// take from current datetime

	opt.Modified = timeNow.String() 		// take from current datetime

	// generate ID of new comment
	hasher := ripemd160.New()
	str_mix := fmt.Sprint("%s | %s | %s | %s", opt.Author, opt.PostID, opt.Date, service.RandSeq(8))
	hasher.Write([]byte(str_mix))
	buf_id := hasher.Sum(nil)
	opt.ID = strings.ToUpper(hex.EncodeToString(buf_id))

	//
	opt_arr, err := json.Marshal(opt)
	if (err != nil) {
		render(w, "comment_create", ActionResult{ Status: "error", Message: err.Error(), Data: opt })
		return
	}
	opt_str := strings.ToUpper(hex.EncodeToString(opt_arr))

	// sign the transaction
	var private_key crypto.PrivKeyEd25519
	private_key, ok = session.Values["private_key"].(crypto.PrivKeyEd25519)
	if (!ok) {
		log.Fatal(err.Error())
		return
	}
	sign := private_key.Sign([]byte(opt_str))
	sign_str := strings.ToUpper(hex.EncodeToString(sign.Bytes()))
	sign_str = sign_str[2:len(sign_str)]

	tx := app.OperationEnvelope {
		Type: "CommentCreateOperation",
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

	tx_json := string(byte_arr[:])
	log.Println("tx_json=", tx_json)

	tx_json_hex := make([]byte, len(tx_json) * 2)
	hex.Encode(tx_json_hex, []byte(tx_json))
	log.Println("tx_json_hex", string(tx_json_hex[:]))

	service.TM_broadcast_tx_commit(string(tx_json_hex[:]))

	// delay sometime then Redirect to the  post
	time.Sleep(1000 * time.Millisecond) // 1s
	http.Redirect(w, req, "/post?p=" + opt.PostID  , http.StatusFound)
}
