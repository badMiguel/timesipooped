package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"timesipooped.fyi/internal"

	// todo remove on production
	"github.com/joho/godotenv"
)

var (
	userInfo = make(map[string]*internal.UserInfo)
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

	authConf := internal.NewOAuthConfig()
	http.HandleFunc("/login", internal.HandleLogin(authConf))
	http.HandleFunc("/login/callback", internal.HandleCallback(authConf, userInfo))

	internal.ConnectDB()

	http.HandleFunc("/poop/add", internal.AddPoop)

	http.HandleFunc("/", test)

	log.Fatal(http.ListenAndServe(os.Getenv("SERVER_PORT"), nil))
}
