package main

import (
	"log"
	"net/http"
	"os"

	"timesipooped.fyi/internal"

	// todo remove on production
	"github.com/joho/godotenv"
)

var (
	userInfo = make(map[string]*internal.UserInfo)
)

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		panic(err)
	}

	// var userInfo *internal.UserInfo
	authConf := internal.NewOAuthConfig()
	http.HandleFunc("/login", internal.HandleLogin(authConf))
	http.HandleFunc("/login/callback", internal.HandleCallback(authConf, userInfo))

	log.Fatal(http.ListenAndServe(os.Getenv("SERVER_PORT"), nil))
}
