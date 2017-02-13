package web

import (
	"net/http"
	"github.com/dawn-network/glogchain/service"
	"log"
	"encoding/json"
	"reflect"
	"github.com/dawn-network/glogchain/app"
	"github.com/dawn-network/glogchain/gopressdb"
)

func Steem_GetPost_Handler(w http.ResponseWriter, req *http.Request) {
	if (req.Method != "POST") {
		render(w, "steem_getpost", ActionResult{Status: "success", Message: "ok", Data: map[string]interface{}{ }})
		return
	}

	// handling POST here
	steem_post_link := req.FormValue("steem_post_link")
	log.Println("steem_post_link", steem_post_link)

	str_json_response, err := service.Steem_get_content(steem_post_link)
	if (err != nil) {
		render(w, "steem_getpost", ActionResult{Status: "error", Message:err.Error(), Data: map[string]interface{}{ }})
		return
	}

	//log.Println("str_json_response", str_json_response)
	/**
	{
	   "id":1,
	   "result":{
	      "id":1174576,
	      "author":"baabeetaa",
	      "permlink":"open-source-steem-blog-coming-soon",
	      "category":"beyondbitcoin",
	      "parent_author":"",
	      "parent_permlink":"beyondbitcoin",
	      "title":"Open Source Steem Blog: Coming Soon!",
	      "body":"<html>\n<p><strong>Steem + Wordpress Realtime Integration</strong></p>\n<p><a href=\"https://postimg.org/image/67y9g8j1b/\">Home page</a>:</p>\n<p><img src=\"https://s19.postimg.org/6kpnmf1b7/home.png\" width=\"558\" height=\"800\"/></p>\n<p><a href=\"https://postimg.org/image/h69iyf7mn/\">Category page</a>:</p>\n<p><img src=\"https://s19.postimg.org/wezgc71b7/category.png\" width=\"558\" height=\"800\"/></p>\n<p><a href=\"https://postimg.org/image/m7gwzsf33/\">Post page</a>:</p>\n<p><img src=\"https://s19.postimg.org/ksfcb2e03/post.png\" width=\"558\" height=\"800\"/></p>\n<p>Thanks @officialfuzzy and @faddat for helping me get getting started with this project.</p>\n</html>",
	      "json_metadata":"{\"tags\":[\"beyondbitcoin\",\"steemit\",\"technology\",\"programming\"],\"users\":[\"officialfuzzy\",\"faddat\"],\"image\":[\"https://s19.postimg.org/6kpnmf1b7/home.png\",\"https://s19.postimg.org/wezgc71b7/category.png\",\"https://s19.postimg.org/ksfcb2e03/post.png\"],\"links\":[\"https://postimg.org/image/67y9g8j1b/\",\"https://postimg.org/image/h69iyf7mn/\",\"https://postimg.org/image/m7gwzsf33/\"]}",
	      "last_update":"2016-10-14T03:21:36",
	      "created":"2016-10-14T03:21:36",
	      "active":"2016-10-22T22:46:06",
	      "last_payout":"2016-11-14T19:52:42",
	      "depth":0,
	      "children":19,
	      "children_rshares2":"0",
	      "net_rshares":0,
	      "abs_rshares":0,
	      "vote_rshares":0,
	      "children_abs_rshares":0,
	      "cashout_time":"1969-12-31T23:59:59",
	      "max_cashout_time":"1969-12-31T23:59:59",
	      "total_vote_weight":0,
	      "reward_weight":10000,
	      "total_payout_value":"176.845 SBD",
	      "curator_payout_value":"58.908 SBD",
	      "author_rewards":666346,
	      "net_votes":211,
	      "root_comment":1174576,
	      "mode":"archived",
	      "max_accepted_payout":"1000000.000 SBD",
	      "percent_steem_dollars":10000,
	      "allow_replies":true,
	      "allow_votes":true,
	      "allow_curation_rewards":true,
	      "url":"/beyondbitcoin/@baabeetaa/open-source-steem-blog-coming-soon",
	      "root_title":"Open Source Steem Blog: Coming Soon!",
	      "pending_payout_value":"0.000 SBD",
	      "total_pending_payout_value":"0.000 SBD",
	      "active_votes":[

	      ],
	      "replies":[

	      ],
	      "author_reputation":"1719768882816",
	      "promoted":"0.000 SBD"
	   }
	}
	 */

	//var objmap map[string]*json.RawMessage
	var res map[string]interface{}
	err = json.Unmarshal([]byte(str_json_response), &res)
	if (err != nil) {
		render(w, "steem_getpost", ActionResult{Status: "error", Message:err.Error(), Data: map[string]interface{}{ }})
		return
	}

	//log.Println("objmap", objmap)
	res_result_i := res["result"]
	log.Println("TypeOf res_result", reflect.TypeOf(res_result_i).String())
	if (reflect.TypeOf(res_result_i).String() != "map[string]interface {}") {
		render(w, "steem_getpost", ActionResult{Status: "error", Message:err.Error(), Data: map[string]interface{}{ }})
		return
	}
	res_result := res_result_i.(map[string]interface {})

	res_result_title := res_result["title"].(string)
	res_result_body := res_result["body"].(string)
	json_metadata_str := res_result["json_metadata"]

	log.Println("res_result_title", res_result_title)
	log.Println("res_result_body", res_result_body)
	log.Println("json_metadata_str", json_metadata_str)

	var json_metadata map[string]interface{}
	err = json.Unmarshal([]byte(json_metadata_str.(string)), &json_metadata)
	if (err != nil) {
		render(w, "steem_getpost", ActionResult{Status: "error", Message:err.Error(), Data: map[string]interface{}{ }})
		return
	}

	//tags_i := json_metadata["tags"].([]interface {})
	//log.Println("tags", tags_i)
	//log.Println("TypeOf tags", reflect.TypeOf(tags_i).String())
	//if (reflect.TypeOf(tags_i).String() != "[]interface {}") {
	//	render(w, "steem_getpost", ActionResult{Status: "error", Message:err.Error(), Data: map[string]interface{}{ }})
	//	return
	//}

	//tags := make([]string, len(tags_i))
	//for i, v := range tags_i {
	//	tags[i] = v.(string)
	//}

	tags, err := json.Marshal(json_metadata["tags"])
	if (err != nil) {
		render(w, "steem_getpost", ActionResult{Status: "error", Message:err.Error(), Data: map[string]interface{}{ }})
		return
	}
	tag_str := string(tags[:])
	log.Println("tags", tag_str)

	//
	images_i := json_metadata["image"].([]interface {})
	if (reflect.TypeOf(images_i).String() != "[]interface {}") {
		render(w, "steem_getpost", ActionResult{Status: "error", Message:err.Error(), Data: map[string]interface{}{ }})
		return
	}
	images := make([]string, len(images_i))
	for i, v := range images_i {
		images[i] = v.(string)
	}

	featured_image := ""
	if (len(images) > 0) {
		featured_image = images[0]
	}

	/////////////////
	var opt app.PostCreateOperation
	//opt.ID = "" // hex string generate by ripemd160 hash the post, will be assinged latter
	//opt.PostAuthor = user.ID 				// take from session
	//opt.PostDate = timeNow.String() 			// take from current datetime
	opt.PostContent = res_result_body
	opt.PostTitle = res_result_title
	//opt.PostModified = timeNow.String() 		// take from current datetime
	opt.Thumb = featured_image
	opt.Cat = tag_str 				// Categories in json string array

	render(w, "PostCreate", ActionResult{Status: "success", Message: "ok", Data: opt })

	//render(w, "steem_getpost", ActionResult{Status: "success", Message: "ok", Data: map[string]interface{}{
	//	"str_json_response": str_json_response,
	//}})
}


func Steem_Posting_Handler(w http.ResponseWriter, req *http.Request) {
	p := req.FormValue("p")
	post := db.GetPost(p)

	render(w, "steem_posting", ActionResult{Status: "success", Message: "ok", Data: post })
}