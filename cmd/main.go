package main

import (
	"database/sql"
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

func connectDB() bool {
	db, err := sql.Open("sqlite3", "../internal/poop.db")
	if err != nil {
		log.Printf("Error opening/creating sqlite3 file: %v\n", err)
		return false
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Printf("Error connecting to the sqlite3 file: %v\n", err)
		return false
	}

	return true
}

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		panic(err)
	}

	for i := 0; i < 3; {
		if success := connectDB(); success {
			break
		}
		i++
		log.Printf("Attempt %d failed. Trying again...\n", i)
	}

	authConf := internal.NewOAuthConfig()
	http.HandleFunc("/login", internal.HandleLogin(authConf))
	http.HandleFunc("/login/callback", internal.HandleCallback(authConf, userInfo))

	http.HandleFunc("/poop/add", internal.AddPoop)

	log.Fatal(http.ListenAndServe(os.Getenv("SERVER_PORT"), nil))
}
