package main

import (
	"database/sql"
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

func connectDB() error {
	db, err := sql.Open("sqlite3", "../internal/poop.db")
	if err != nil {
		return fmt.Errorf("Error opening/creating sqlite3 file: %v\n", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return fmt.Errorf("Error connecting to the sqlite3 file: %v\n", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
        userId INTEGER PRIMARY KEY,
        email TEXT UNIQUE,
        givenName TEXT, 
        familyName TEXT, 
        joined DATETIME DEFAULT CURRENT_TIMESTAMP,
        poopTotal INTEGER DEFAULT 0,
        failedTotal INTEGER DEFAULT 0
    )`)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("Error creating new users table: %v\n", err))
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS poop (
        poopId INTEGER PRIMARY KEY AUTOINCREMENT,
        date DATETIME DEFAULT CURRENT_TIMESTAMP,
        userId INTEGER,
        FOREIGN KEY(userId) REFERENCES users(userId) ON DELETE CASCADE
    )`)
	if err != nil {
		return fmt.Errorf("Error creating new poop table: %v\n", err)
	}

	return nil
}

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

	for i := 0; i < 10; {
		if err := connectDB(); err == nil {
			log.Println("Successfully connected to the database")
			break
		} else {
			log.Println(err)
		}

		i++
		log.Printf("Attempt %d failed. Trying again...\n", i)

		if i == 3 {
			panic("Maximum attempts reached. Could not connect to DB")
		}
	}

	authConf := internal.NewOAuthConfig()
	http.HandleFunc("/login", internal.HandleLogin(authConf))
	http.HandleFunc("/login/callback", internal.HandleCallback(authConf, userInfo))

	http.HandleFunc("/poop/add", internal.AddPoop)

	http.HandleFunc("/", test)

	log.Fatal(http.ListenAndServe(os.Getenv("SERVER_PORT"), nil))
}
