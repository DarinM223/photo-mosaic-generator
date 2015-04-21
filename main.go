package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const PORT = "8080"

type ConfigJson struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type AccessTokenResponse struct {
	AccessToken string          `json:"access_token"`
	User        AccessTokenUser `json:"user"`
}

type AccessTokenUser struct {
	Id             string `json:"id"`
	Username       string `json:"username"`
	FullName       string `json:"full_name"`
	ProfilePicture string `json:"profile_picture"`
}

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

func RetrieveAccessToken(code string, config ConfigJson) (AccessTokenResponse, error) {
	/// request access token
	reqForm := url.Values{
		"client_id":     {config.ClientId},
		"client_secret": {config.ClientSecret},
		"grant_type":    {"authorization_code"},
		"redirect_uri":  {"http://localhost:" + PORT},
		"code":          {code},
	}

	resp, err := http.PostForm("https://api.instagram.com/oauth/access_token", reqForm)

	if err != nil {
		return AccessTokenResponse{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return AccessTokenResponse{}, err
	}

	/// parse response json into struct
	var parsedResponse AccessTokenResponse
	err = json.Unmarshal(body, &parsedResponse)

	if err != nil {
		return AccessTokenResponse{}, err
	}

	return parsedResponse, nil
}

func main() {
	config, _ := ParseConfig("./config.json")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		RootHandler(w, r, config)
	})
	http.HandleFunc("/instagram_authorize", func(w http.ResponseWriter, r *http.Request) {
		InstagramLoginRedirect(w, r, config)
	})

	fmt.Println("Listening on port", PORT)
	http.ListenAndServe(":"+PORT, nil)
}

/// Request handlers

func RootHandler(w http.ResponseWriter, r *http.Request, config ConfigJson) {
	code := r.URL.Query().Get("code")

	if len(code) == 0 {
		html := `
		<html>
		<head>
		<title>Home page</title>
		</head>
		<body>
		Home page!
		Log in with <a href="/instagram_authorize">Instagram</a>
		</body>
		</html>
		`
		w.Write([]byte(html))
	} else {
		parsedResponse, err := RetrieveAccessToken(code, config)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println("Access Token:", parsedResponse.AccessToken)
	}
}

func InstagramLoginRedirect(w http.ResponseWriter, r *http.Request, config ConfigJson) {
	redirectUri := url.QueryEscape("http://localhost:" + PORT)

	responseType := "code"

	http.Redirect(w, r,
		"https://api.instagram.com/oauth/authorize/?client_id="+config.ClientId+"&redirect_uri="+
			redirectUri+"&response_type="+responseType,
		301)
}
