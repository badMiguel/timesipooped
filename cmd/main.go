package main

import (
	"database/sql"
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
	userInfo = make(map[string]*database.UserInfo)
)

func test(w http.ResponseWriter, r *http.Request) {
	for i := range userInfo {
		fmt.Fprintf(w, "%v\n", *userInfo[i])
	}
}

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		panic(err)
	}

	db, err := sql.Open("sqlite3", "../internal/database/db.sqlite3")
	if err != nil {
		panic(fmt.Sprintf("Error opening/creating sqlite3 file: %v\n", err))
	}
	defer db.Close()
	database.InitializeDB(db)

	authConf := login.NewOAuthConfig()
	http.HandleFunc("/login", login.HandleLogin(authConf))
	http.HandleFunc("/login/callback", login.HandleCallback(authConf, userInfo, db))

	http.HandleFunc("/poop/add", poop.AddPoop)

	// REMOVE AFTER
	http.HandleFunc("/", test)
	log.Fatal(http.ListenAndServe(os.Getenv("SERVER_PORT"), nil))
}
