package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func initialize() error {
	db, err := sql.Open("sqlite3", "./database/db.sqlite3")
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
	log.Println("Created users table")

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS poop (
        poopId INTEGER PRIMARY KEY AUTOINCREMENT,
        date DATETIME DEFAULT CURRENT_TIMESTAMP,
        userId INTEGER,
        FOREIGN KEY(userId) REFERENCES users(userId) ON DELETE CASCADE
    )`)
	if err != nil {
		return fmt.Errorf("Error creating new poop table: %v\n", err)
	}
	log.Println("Created poop table")

	return nil
}

func StartDB() {
	for i := 0; i < 10; {
		if err := initialize(); err == nil {
			log.Println("Successfully connected to the database")
			break
		} else {
			log.Println(err)
		}

		i++
		log.Printf("Attempt %d failed. Trying again...\n", i)

		if i == 10 {
			panic("Maximum attempts reached. Could not connect to DB")
		}
	}
}
