package web

import (
	"encoding/json"
	"golang.org/x/crypto/ripemd160"
	"fmt"
	"github.com/dawn-network/glogchain/service"
	"strings"
	"encoding/hex"
	"log"
	"time"
	"net/http"
	"github.com/dawn-network/glogchain/db"
	"github.com/tendermint/go-crypto"
	"github.com/dawn-network/glogchain/app"
)

func PostCreateHandler(w http.ResponseWriter, req *http.Request) {
	log.Println("---PostCreateHandler---------------------------------------------------------------------")

	var opt app.PostCreateOperation

	if req.Method != "POST" {
		log.Println("PostCreateHandler GET")
		context := Context{Title: "PostCreate", Data: opt}
		render(w, "PostCreate", context)
		return
	}

	log.Println("PostCreateHandler POST")

	session := GetSession(req)
	user, ok := session.Values["user"].(db.User)
	if !ok {
		panic("PostCreateHandler: wtf?? session user is nil")
		return
	}

	timeNow := time.Now()

	/////////
	opt.ID = "" // hex string generate by ripemd160 hash the post, will be assinged latter
	opt.PostAuthor = user.ID 				// take from session
	opt.PostDate = timeNow.String() 			// take from current datetime
	opt.PostContent = req.FormValue("PostContent")
	opt.PostTitle = req.FormValue("Title")
	opt.PostModified = timeNow.String() 		// take from current datetime
	opt.Thumb = req.FormValue("Thumb")
	opt.Cat = req.FormValue("Categories") 	// Categories in json string array

	// very basic from validating
	if ((len(opt.PostTitle) < 6) || (len(opt.PostContent) < 6) || (len(opt.Thumb) < 6)) {
		render(w, "PostCreate", ActionResult{Status: "error", Message: "field must be at least 6 characters", Data: opt})
		return
	}

	/////////////////////////////////////////
	// Save the post content to IPFS and return its hash to opt.PostContent
	mhash, err := service.Ipfs_add([]byte(opt.PostContent))
	if (err != nil) {
		render(w, "PostCreate", ActionResult{Status: "error", Message: "Can not save content to IPFS", Data: opt})
		return
	}
	opt.PostContent = mhash // set the body to hash because content already stored in IPFS

	// validating Categories
	opt.Cat = strings.ToLower(opt.Cat)
	if (len(opt.Cat) < 6) {
		opt.Cat = `[]`
	}

	cats_string := []string{}
	err = json.Unmarshal([]byte(opt.Cat), &cats_string)
	if (err != nil) {
		render(w, "PostCreate", ActionResult{Status: "error", Message: "Categories json array string is invalid", Data: opt})
		return
	}

	// generate id of new post
	hasher := ripemd160.New()
	str_mix := fmt.Sprint("%s | %s | %s | %s", opt.PostAuthor, opt.PostTitle, opt.PostDate, service.RandSeq(8))
	hasher.Write([]byte(str_mix))
	buf_id := hasher.Sum(nil)
	opt.ID = strings.ToUpper(hex.EncodeToString(buf_id))

	log.Println("PostCreateHandler", "id=", opt.ID, "PostAuthor=", opt.PostAuthor, "PostDate=",
		opt.PostDate, "Title=", opt.PostTitle, "PostModified=", opt.PostModified, "Thumb=",
		opt.Thumb, "Categories=", opt.Cat, "PostContent=", opt.PostContent)

	opt_arr, err := json.Marshal(opt)
	if (err != nil) {
		render(w, "PostCreate", ActionResult{Status: "error", Message: err.Error(), Data: opt })
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
		Type: "PostCreateOperation",
		Operation: opt_str,
		Signature: sign_str,
		Pubkey: private_key.PubKey().KeyString(),
		Fee: 0,
	}

	byte_arr, err := json.Marshal(tx)
	if err != nil {
		log.Fatal("PostCreateHandler Cannot encode to JSON ", err)
		return
	}

	tx_json := string(byte_arr[:])
	log.Println("PostCreateHandler tx_json=", tx_json)

	tx_json_hex := make([]byte, len(tx_json) * 2)
	hex.Encode(tx_json_hex, []byte(tx_json))
	log.Println("PostCreateHandler tx_json_hex", string(tx_json_hex[:]))

	service.TM_broadcast_tx_commit(string(tx_json_hex[:]))

	// delay sometime then Redirect to the new post
	time.Sleep(1000 * time.Millisecond) // 1s
	http.Redirect(w, req, "/post?p=" + opt.ID  , http.StatusFound)
}


func PostEditHandler(w http.ResponseWriter, req *http.Request) {
	log.Println("---PostEditHandler---------------------------------------------------------------------")

	p := req.FormValue("p")
	post, err := db.GetPost(p)

	session := GetSession(req)
	user, ok := session.Values["user"].(db.User)
	if (!ok) {
		panic("PostCreateHandler: wtf?? session user is nil")
		return
	}

	// If method is GET serve an html
	if req.Method != "POST" {
		log.Println("PostEditHandler GET")

		// TODO: Only the author can edit the post
		// https://github.com/dawn-network/glogchain/issues/8

		// check the current logged user is the author of the post or not
		if (user.ID != post.PostAuthor) {
			render(w, "post_edit",
				ActionResult{
					Status: "error",
					Message: "Only the author can edit the post",
					Data: map[string]interface{}{ },
				})
			return
		}

		context := Context{Title: "PostEdit", Data: post}
		render(w, "post_edit", context)
		return
	}

	log.Println("PostEditHandler POST")
	timeNow := time.Now()

	/////////
	post.ID = p
	post.PostAuthor = user.ID 				// take from session
	//post.PostDate = timeNow.String() 			// take from current datetime
	post.PostContent = req.FormValue("PostContent")
	post.PostTitle = req.FormValue("Title")
	post.PostModified = timeNow.String() 		// take from current datetime
	post.Thumb = req.FormValue("Thumb")
	post.Cat = req.FormValue("Categories") 	// Categories in json string array

	// very basic from validating
	if ((len(post.PostTitle) < 6) || (len(post.PostContent) < 6) || (len(post.Thumb) < 6)) {
		render(w, "post_edit",
			ActionResult{
				Status: "error",
				Message: "field must be at least 6 characters",
				Data: post,
			})
		return
	}

	// validating Categories
	post.Cat = strings.ToLower(post.Cat)
	if (len(post.Cat) < 6) {
		post.Cat = `[]`
	}

	cats_string := []string{}
	err = json.Unmarshal([]byte(post.Cat), &cats_string)
	if (err != nil) {
		render(w, "post_edit",
			ActionResult{
				Status: "error",
				Message: "Categories json array string is invalid",
				Data: post,
			})
		return
	}


	log.Println("PostEditHandler", "id=", post.ID, "PostAuthor=", post.PostAuthor, "PostDate=",
		post.PostDate, "Title=", post.PostTitle, "PostModified=", post.PostModified, "Thumb=",
		post.Thumb, "Categories=", post.Cat, "PostContent=", post.PostContent)

	opt_arr, err := json.Marshal(post)
	if (err != nil) {
		render(w, "post_edit",
			ActionResult{
				Status: "error",
				Message: err.Error(),
				Data: post,
			})
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
		Type: "PostEditOperation",
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

	tx_json := string(byte_arr[:])
	log.Println("PostEditHandler tx_json=", tx_json)

	tx_json_hex := make([]byte, len(tx_json) * 2)
	hex.Encode(tx_json_hex, []byte(tx_json))
	log.Println("PostEditHandler tx_json_hex", string(tx_json_hex[:]))

	service.TM_broadcast_tx_commit(string(tx_json_hex[:]))

	// delay sometime then Redirect to the new post
	time.Sleep(1000 * time.Millisecond) // 1s
	http.Redirect(w, req, "/post?p=" + post.ID  , http.StatusFound)
}
