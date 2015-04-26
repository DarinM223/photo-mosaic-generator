package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
)

/// json retrieval functions

func RetrieveAccessToken(code string, config ConfigJson) (AccessTokenResponse, error) {
	// request access token
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

	// parse response json into struct
	var parsedResponse AccessTokenResponse
	err = json.Unmarshal(body, &parsedResponse)
	if err != nil {
		return AccessTokenResponse{}, err
	}

	return parsedResponse, nil
}

func RetrieveImageData(tag string, access_token string) (TagSearchResponse, error) {
	resp, err := http.Get("https://api.instagram.com/v1/tags/" + url.QueryEscape(tag) +
		"/media/recent?access_token=" + url.QueryEscape(access_token))
	if err != nil {
		return TagSearchResponse{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return TagSearchResponse{}, err
	}

	var tagSearchResponse TagSearchResponse
	err = json.Unmarshal(body, &tagSearchResponse)
	if err != nil {
		return TagSearchResponse{}, err
	}

	return tagSearchResponse, nil
}

/// Request handlers

func RootHandler(w http.ResponseWriter, r *http.Request, config ConfigJson) {
	code := r.URL.Query().Get("code")
	cookie, cookie_err := r.Cookie(AUTH_COOKIE_NAME)
	_ = cookie

	if len(code) == 0 && cookie_err != nil { // if auth code sent
		t, _ := template.ParseFiles("templates/login.html")
		t.Execute(w, nil)
	} else if cookie_err != nil { // if cookie not set yet
		parsedResponse, err := RetrieveAccessToken(code, config)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:   AUTH_COOKIE_NAME,
			Value:  parsedResponse.AccessToken,
			MaxAge: 0,
		})

		http.SetCookie(w, &http.Cookie{
			Name:   USERNAME_COOKIE_NAME,
			Value:  parsedResponse.User.FullName,
			MaxAge: 0,
		})

		// create template
		page := struct {
			Name string
		}{
			parsedResponse.User.FullName,
		}

		t, _ := template.ParseFiles("templates/index.html")
		t.Execute(w, page)
	} else { // if authentication token already set in cookie
		username_cookie, err := r.Cookie(USERNAME_COOKIE_NAME)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		// create template
		page := struct {
			Name string
		}{
			username_cookie.Value,
		}

		t, _ := template.ParseFiles("templates/index.html")
		t.Execute(w, page)
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

func InstagramSearch(w http.ResponseWriter, r *http.Request, config ConfigJson) {
	tag_query := r.URL.Query().Get("query")
	auth_cookie, err := r.Cookie(AUTH_COOKIE_NAME)
	if err != nil {
		http.Error(w, "You have to be authorized to access this page", http.StatusUnauthorized)
		return
	}
	access_token := auth_cookie.Value

	tagSearchResponse, err := RetrieveImageData(tag_query, access_token)
	_ = tagSearchResponse
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
