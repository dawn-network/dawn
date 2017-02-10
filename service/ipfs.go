package service

import (
	"io"
	"bytes"
	"mime/multipart"
	"net/http"
	"io/ioutil"
	"log"
	"github.com/baabeetaa/glogchain/app"
	"encoding/json"
	"github.com/pkg/errors"
)

// https://ipfs.io/docs/api/#apiv0add
func Ipfs_add(r io.Reader) (mhash string, err error) {

	//////////////////////////////////
	// DONT USE THE GO IPFS API CLIENT, IT BREAKS GLIDE
	//log.Println("Ipfs_add")
	//
	//s := shell.NewShell(app.GlogchainConfigGlobal.IpFsAPI)
	//
	//mhash, err = s.Add(r)
	//
	//log.Println("mhash", mhash)
	//////////////////////////////////////////

	// Prepare a form that you will submit to that URL.
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	fw, err := w.CreateFormFile("file", "myfile")
	if err != nil {
		log.Println(err.Error())
		return
	}
	if _, err = io.Copy(fw, r); err != nil {
		log.Println(err.Error())
		return
	}

	//// Add the other fields
	//if fw, err = w.CreateFormField("key"); err != nil {
	//	return
	//}
	//if _, err = fw.Write([]byte("KEY")); err != nil {
	//	return
	//}
	// Don't forget to close the multipart writer.
	// If you don't close it, your request will be missing the terminating boundary.
	w.Close()

	// Now that you have a form, you can submit it to your handler.
	req, err := http.NewRequest("POST", app.GlogchainConfigGlobal.IpFsAPI + "/api/v0/add", &b)
	if err != nil {
		log.Println(err.Error())
		return
	}
	// Don't forget to set the content type, this will contain the boundary.
	req.Header.Set("Content-Type", w.FormDataContentType())

	// Submit the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err.Error())
		return
	}

	//// Check the response
	//if resp.StatusCode != http.StatusOK {
	//	err = fmt.Errorf("bad status: %s", resp.Status)
	//}

	defer resp.Body.Close()

	log.Println("response Status:", resp.Status)
	log.Println("response Headers:", resp.Header)
	body, err := ioutil.ReadAll(resp.Body)
	if (err != nil) {
		log.Println(err.Error())
		return
	}

	str_json_response := string(body)
	log.Println("str_json_response", str_json_response)
	/**
	Neet to parse the json to get the hash value
	{"Name":"myfile","Hash":"QmZU9Ln7eoKBc5nL9LyatzhxL7ce4yST7oXxT6N5V5X6LZ"}
	 */
	var jobj map[string]interface{}
	err = json.Unmarshal([]byte(str_json_response), &jobj)
	if (err != nil) {
		log.Println(err.Error())
		return
	}

	mhash, ok := jobj["Hash"].(string)
	if !ok {
		err = errors.New("Can not get Hash")
		return
	}

	return
}