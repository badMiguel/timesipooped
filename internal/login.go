package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
)

type OAuthConfig struct {
	ClientId     string
	ClientSecret string
	RedirectURI  string
	OauthURL     string
	TokenURL     string
	UserInfoURL  string
}

type UserInfo struct {
	Email          string `json:"email"`
	Family_name    string `json:"family_name"`
	Given_name     string `json:"given_name"`
	Id             string `json:"id"`
	Name           string `json:"name"`
	Picture        string `json:"picture"`
	Verified_email bool   `json:"verified_email"`
}

func NewOAuthConfig() *OAuthConfig {
	config := &OAuthConfig{}
	config.ClientId = os.Getenv("CLIENT_ID")
	config.ClientSecret = os.Getenv("CLIENT_SECRET")
	config.RedirectURI = os.Getenv("REDIRECT_URI")
	config.OauthURL = "https://accounts.google.com/o/oauth2/auth"
	config.TokenURL = "https://accounts.google.com/o/oauth2/token"
	config.UserInfoURL = "https://www.googleapis.com/oauth2/v1/userinfo?alt=json"

	return config
}

func HandleLogin(authConf *OAuthConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authUrl := fmt.Sprintf(
			"%s?response_type=code&client_id=%s&redirect_uri=%s&scope=openid%%20email%%20https://www.googleapis.com/auth/userinfo.profile",
			authConf.OauthURL,
			authConf.ClientId,
			url.QueryEscape(authConf.RedirectURI),
		)
		http.Redirect(w, r, authUrl, http.StatusFound)
	}
}

func HandleCallback(authConf *OAuthConfig, userInfo *UserInfo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code == "" {
			http.Error(w, "Missing code parameter", http.StatusBadRequest)
			return
		}

		resp, err := http.PostForm(authConf.TokenURL, url.Values{
			"code":          {code},
			"client_id":     {authConf.ClientId},
			"client_secret": {authConf.ClientSecret},
			"redirect_uri":  {authConf.RedirectURI},
			"grant_type":    {"authorization_code"},
		})
		if err != nil {
			log.Printf("Error exchanging code for token: %v\n", err)
			http.Error(w, "Failed to exchange code", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		var tokenData map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&tokenData); err != nil {
			http.Error(w, "Failed to decode token response", http.StatusInternalServerError)
			return
		}

		accessToken := tokenData["access_token"].(string)

		req, err := http.NewRequest("GET", authConf.UserInfoURL, nil)
		if err != nil {
			log.Printf("Error creating user info request: %v\n", err)
			http.Error(w, "Failed to create user info request", http.StatusInternalServerError)
		}
		req.Header.Set("Authorization", "Bearer "+accessToken)
		client := &http.Client{}
		resp, err = client.Do(req)
		if err != nil {
			log.Printf("Error fetching user info: %v\n", err)
			http.Error(w, "Failed to fetch user info", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
			log.Printf("Error decoding user info: %v\n", err)
			http.Error(w, "Failed to decode user info", http.StatusInternalServerError)
			return
		}
	}
}
