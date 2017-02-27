package service

import (
	"net/url"
	"log"
	"strings"
	"errors"
	"fmt"
	"bytes"
	"io/ioutil"
	"net/http"
)

const (
	// https://steem.yt, http://steem.yt:8090, http://138.201.198.175:8090 https://node.steem.ws/ wss://node.steem.ws/
	Steem_Rpc_Url = "https://node.steem.ws/"
)


/**
 Get a post from steem by given post URL. eg.,
 https://steemit.com/beyondbitcoin/@faddat/application-design-blockchain-account-index

 -------
 https://github.com/go-steem/rpc <== usefull ?! OR maybe dont need it
 Steem API: http://steemroll.io/api-explorer/#method=get_content&params=["lantto","steem-api-explorer"]
 */
func Steem_get_content(steemit_link string) (str_json_response string, err error) {
	/////////////////////////////
	// First, we have to parse the steemit_url to get the author and the post permlink
	url_steem, err := url.Parse(steemit_link)
	if (err != nil) {
		log.Println(err.Error())
		return
	}

	log.Println("url_steem", url_steem)

	//log.Println("url_steem.RawPath", url_steem.RawPath)
	//log.Println("url_steem.Path", url_steem.Path)

	// RawPath should be: /beyondbitcoin/@faddat/application-design-blockchain-account-index
	strs := strings.Split(url_steem.Path, "/")
	log.Println("strs", strs)

	// find the path content of user (@faddat)
	idx := -1
	for i := 0; i < len(strs); i++ {
		if (strings.HasPrefix(strs[i], "@")) {
			idx = i		// hold the index
			break		// break the loop
		}
	}

	if (idx < 0) {
		// invalid steemit_link
		err = errors.New("Invalid URL!!!!")
		return
	}

	// remove the '@' from author username
	steem_post_author := strs[idx]
	steem_post_author = steem_post_author[1:len(steem_post_author)]
	log.Println("steem_post_author", steem_post_author)

	steem_post_permlink := strs[idx + 1] 		// hold the permlink
	log.Println("steem_post_permlink", steem_post_permlink)


	/////////////////////////////
	// OK, now we have the author and permlink, let get the post content from steem api

	// build the post content
	json_str := fmt.Sprintf(`{"id": 1, "jsonrpc": "2.0", "method": "call", "params": ["database_api", "get_content", ["%s", "%s"]] }`, steem_post_author, steem_post_permlink)
	log.Println("json_str", json_str)

	// make the post request
	req, err := http.NewRequest("POST", Steem_Rpc_Url, bytes.NewBuffer([]byte(json_str)))
	if (err != nil) {
		log.Println(err.Error())
		return
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if (err != nil) {
		log.Println(err.Error())
		return
	}

	defer resp.Body.Close()

	log.Println("response Status:", resp.Status)
	log.Println("response Headers:", resp.Header)
	body, err := ioutil.ReadAll(resp.Body)
	if (err != nil) {
		log.Println(err.Error())
		return
	}

	str_json_response = string(body)
	//log.Println("response Body:", str_json_response)

	return
}