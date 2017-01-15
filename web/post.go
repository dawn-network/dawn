package web

import (
	"encoding/json"
	"golang.org/x/crypto/ripemd160"
	"fmt"
	"github.com/baabeetaa/glogchain/service"
	"strings"
	"encoding/hex"
	"log"
	"github.com/baabeetaa/glogchain/protocol"
	"time"
	"net/http"
	"github.com/baabeetaa/glogchain/db"
)

func PostCreateHandler(w http.ResponseWriter, req *http.Request) {
	log.Println("---PostCreateHandler---------------------------------------------------------------------")

	var post db.Post

	// If method is GET serve an html login page
	if req.Method != "POST" {
		log.Println("PostCreateHandler GET")
		context := Context{Title: "PostCreate", Data: post}
		render(w, "PostCreate", context)
		return
	}

	log.Println("PostCreateHandler POST")

	session := GetSession(req)
	user, ok := session.Values["user"].(*db.User)
	if !ok {
		panic("PostCreateHandler: wtf?? session user is nil")
		return
	}

	timeNow := time.Now()

	/////////
	post.ID = "" // hex string generate by ripemd160 hash the post, will be assinged latter
	post.PostAuthor = user.ID 				// take from session
	post.PostDate = timeNow.String() 			// take from current datetime
	post.PostContent = req.FormValue("PostContent")
	post.PostTitle = req.FormValue("Title")
	post.PostModified = timeNow.String() 		// take from current datetime
	post.Thumb = req.FormValue("Thumb")
	post.Cat = req.FormValue("Categories") 	// Categories in json string array

	// very basic from validating
	if ((len(post.PostTitle) < 6) || (len(post.PostContent) < 6) || (len(post.Thumb) < 6)) {
		render(w, "PostCreate",
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
	err := json.Unmarshal([]byte(post.Cat), &cats_string)
	if (err != nil) {
		render(w, "PostCreate",
			ActionResult{
				Status: "error",
				Message: "Categories json array string is invalid",
				Data: post,
			})
		return
	}

	// generate id of new post
	hasher := ripemd160.New()
	str_mix := fmt.Sprint("%s | %s | %s | %s", post.PostAuthor, post.PostTitle, post.PostDate, service.RandSeq(8))
	hasher.Write([]byte(str_mix))
	buf_id := hasher.Sum(nil)
	post.ID = strings.ToUpper(hex.EncodeToString(buf_id))

	log.Println("PostCreateHandler", "id=", post.ID, "PostAuthor=", post.PostAuthor, "PostDate=",
		post.PostDate, "Title=", post.PostTitle, "PostModified=", post.PostModified, "Thumb=",
		post.Thumb, "Categories=", post.Cat, "PostContent=", post.PostContent)


	tx := protocol.OperationEnvelope {
		Type: "PostCreateOperation",
		Operation: protocol.PostCreateOperation {
			Fee: 0,
			Post: post },
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
	http.Redirect(w, req, "/post?p=" + post.ID  , http.StatusFound)
}


func PostEditHandler(w http.ResponseWriter, req *http.Request) {
	log.Println("---PostEditHandler---------------------------------------------------------------------")

	p := req.FormValue("p")
	post, err := db.GetPost(p)

	// If method is GET serve an html login page
	if req.Method != "POST" {
		log.Println("PostEditHandler GET")

		context := Context{Title: "PostEdit", Data: post}
		render(w, "post_edit", context)
		return
	}

	log.Println("PostEditHandler POST")

	session := GetSession(req)
	user, ok := session.Values["user"].(*db.User)
	if !ok {
		panic("PostCreateHandler: wtf?? session user is nil")
		return
	}

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


	tx := protocol.OperationEnvelope {
		Type: "PostEditOperation",
		Operation: protocol.PostEditOperation {
			Fee: 0,
			Post: post },
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
