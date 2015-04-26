package main

const PORT = "8080"

const (
	AUTH_COOKIE_NAME     = "access_token_cookie"
	USERNAME_COOKIE_NAME = "username_cookie"
)

type ConfigJson struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

/// Access token response structs
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

/// Tag search response structs
type TagSearchResponse struct {
	Data []TagImage `json:"data"`
}

type TagImage struct {
	Type   string        `json:"type"`
	Images TagImageTypes `json:"images"`
}

type TagImageTypes struct {
	LowRes      TagImageProperties `json:"low_resolution"`
	Thumbnail   TagImageProperties `json:"thumbnail"`
	StandardRes TagImageProperties `json:"standard_resolution"`
}

type TagImageProperties struct {
	Url    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}
