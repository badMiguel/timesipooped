package main

import (
	// "encoding/json"
	"log"
	"net/http"
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

func get(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func post(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func newOAuthConfig() *OAuthConfig {
	config := &OAuthConfig{}
	err := godotenv.Load("../.env")
	if err != nil {
		panic(err)
	}
	config.ClientId = os.Getenv("CLIENT_ID")
	config.ClientSecret = os.Getenv("CLIENT_SECRET")
	config.RedirectURI = "http://localhost:8080/callback"
	config.OauthURL = "https://accounts.google.com/o/oauth2/auth"
	config.TokenURL = "https://accounts.google.com/o/oauth2/token"
	config.UserInfoURL = "https://www.googleapis.com/oauth2/v1/userinfo?alt=json"
	return config
}

func main() {
	http.HandleFunc("/get", get)
	http.HandleFunc("/post", post)

	log.Fatal(http.ListenAndServe(":1234", nil))
}
