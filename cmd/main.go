package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"timesipooped.fyi/internal/database"
	"timesipooped.fyi/internal/auth"
	"timesipooped.fyi/internal/poop"

	_ "github.com/mattn/go-sqlite3"

	// todo remove on production
	"github.com/joho/godotenv"
)

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

	authConf := auth.NewOAuthConfig()
	http.HandleFunc("/auth/login", auth.HandleLogin(authConf))
	http.HandleFunc("/auth/login/callback", auth.HandleCallback(authConf, db))
	http.HandleFunc("/auth/status", auth.HandleStatus(authConf))

	http.HandleFunc("/poop/add", poop.AddPoop)

	log.Fatal(http.ListenAndServe(os.Getenv("SERVER_PORT"), nil))
}
