package main

import (
	// "encoding/json"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	// todo remove on production
	"github.com/joho/godotenv"
)

type OAuthConfig struct {
	ClientId     string
	ClientSecret string
	RedirectURI  string
	OauthURL     string
	TokenURL     string
	UserInfoURL  string
}

type PoopList struct {
	p []string
}

type PoopInfo struct {
	poopTotal  int
	fPoopTotal int

	poopList  *PoopList
	fPoopList *PoopList
}

type User struct {
	Username string    `json:"username"`
	Info     *PoopInfo `json:"info"`
}

func newPoopInfo() *PoopInfo {
	return &PoopInfo{
		poopTotal:  0,
		fPoopTotal: 0,
		poopList:   &PoopList{[]string{}},
		fPoopList:  &PoopList{[]string{}},
	}
}

func newUser(u string) *User {
	return &User{
		Username: u,
		Info:     newPoopInfo(),
	}
}

func handleGetReq(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func handlePostReq(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func newOAuthConfig() *OAuthConfig {
	config := &OAuthConfig{}
	config.ClientId = os.Getenv("CLIENT_ID")
	config.ClientSecret = os.Getenv("CLIENT_SECRET")
	config.RedirectURI = os.Getenv("REDIRECT_URI")
	config.OauthURL = "https://accounts.google.com/o/oauth2/auth"
	config.TokenURL = "https://accounts.google.com/o/oauth2/token"
	config.UserInfoURL = "https://www.googleapis.com/oauth2/v1/userinfo?alt=json"

	return config
}

func handleLogin(authConf *OAuthConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authUrl := fmt.Sprintf(
			"%s?response_type=code&client_id=%s&redirect_uri=%s&scope=openid%%20email",
			authConf.OauthURL,
			authConf.ClientId,
			url.QueryEscape(authConf.RedirectURI),
		)
		http.Redirect(w, r, authUrl, http.StatusFound)
	}
}

func handleCallback(authConf *OAuthConfig) http.HandlerFunc {
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

		var userInfo map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
			log.Printf("Error decoding user info: %v\n", err)
			http.Error(w, "Failed to decode user info", http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "User Info: %v", userInfo)
	}
}

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		panic(err)
	}

	authConf := newOAuthConfig()
	http.HandleFunc("/login", handleLogin(authConf))
	http.HandleFunc("/login/callback", handleCallback(authConf))

	http.HandleFunc("/get", handleGetReq)
	http.HandleFunc("/post", handlePostReq)

	log.Fatal(http.ListenAndServe(os.Getenv("SERVER_PORT"), nil))
}
