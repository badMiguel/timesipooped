package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"timesipooped.fyi/internal/database"
	"timesipooped.fyi/internal/login"
	"timesipooped.fyi/internal/poop"

	_ "github.com/mattn/go-sqlite3"

	// todo remove on production
	"github.com/joho/godotenv"
)

var (
	userInfo = make(map[string]*login.UserInfo)
)

func test(w http.ResponseWriter, r *http.Request) {
	for i := range userInfo {
		fmt.Fprintln(w, *userInfo[i])
	}
}

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		panic(err)
	}

	authConf := login.NewOAuthConfig()
	http.HandleFunc("/login", login.HandleLogin(authConf))
	http.HandleFunc("/login/callback", login.HandleCallback(authConf, userInfo))

	database.StartDB()
	http.HandleFunc("/poop/add", poop.AddPoop)

	// REMOVE AFTER
	http.HandleFunc("/", test)
	log.Fatal(http.ListenAndServe(os.Getenv("SERVER_PORT"), nil))
}
