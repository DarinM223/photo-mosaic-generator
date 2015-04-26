package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func ParseConfig(filename string) (ConfigJson, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return ConfigJson{}, err
	}

	var config ConfigJson
	err = json.Unmarshal(content, &config)
	if err != nil {
		return ConfigJson{}, err
	}

	return config, nil
}

func main() {
	config, _ := ParseConfig("./config.json")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		RootHandler(w, r, config)
	})
	http.HandleFunc("/instagram_authorize", func(w http.ResponseWriter, r *http.Request) {
		InstagramLoginRedirect(w, r, config)
	})
	http.HandleFunc("/instagram_search", func(w http.ResponseWriter, r *http.Request) {
		InstagramSearch(w, r, config)
	})

	fmt.Println("Listening on port", PORT)
	http.ListenAndServe(":"+PORT, nil)
}
