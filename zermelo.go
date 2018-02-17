package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

var baseurl string

// TODO:
// Write functions that do get rqequests to the zermelo api and implement some funcs, maybe a class <- struct zermelo or smth idk

func fetchAuthToken(organisation string, code int) string {
	// Do a http request to get the auth token
	baseurl = "https://" + organisation + ".zportal.nl/api/v3"
	response, err := http.PostForm(baseurl+"/oauth/token", url.Values{"grant_type": {"authorization_code"}, "code": {strconv.Itoa(code)}})
	if err != nil {
		fmt.Println(err)
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return ""
	}

	// Read the body from the response
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
	}

	// Parse json and return token
	var jsonBody map[string]string
	_ = json.Unmarshal(body, &jsonBody) // We expect an error, because the expires_at field is [string]int, we don't need it anyway though
	return jsonBody["access_token"]
}
