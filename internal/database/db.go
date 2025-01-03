package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type UserSchema struct {
	userId      string
	email       string
	givenName   string
	familyName  string
	joined      time.Time
	poopTotal   int
	failedTotal int
}

type UserInfo struct {
	Id            string `json:"id"`
	Email         string `json:"email"`
	FamilyName    string `json:"family_name"`
	GivenName     string `json:"given_name"`
	VerifiedEmail bool   `json:"verified_email"`
}

type PoopSchema struct {
	poopId int
	date   time.Time
	userId string
}

func initialize(db *sql.DB) error {
	if err := db.Ping(); err != nil {
		return fmt.Errorf("Error connecting to the sqlite3 file: %v\n", err)
	}

	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS users (
        userId TEXT PRIMARY KEY,
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
	log.Println("Initialized users table")

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS poop (
        poopId INTEGER PRIMARY KEY AUTOINCREMENT,
        date DATETIME DEFAULT CURRENT_TIMESTAMP,
        userId INTEGER,
        FOREIGN KEY(userId) REFERENCES users(userId) ON DELETE CASCADE
    )`)
	if err != nil {
		return fmt.Errorf("Error creating new poop table: %v\n", err)
	}
	log.Println("Initialized poop table")

	return nil
}

func InitializeDB(db *sql.DB) {
	for i := 0; i < 10; {
		if err := initialize(db); err == nil {
			log.Println("Successfully connected to the database")
			return
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

func addUser(userInfo *UserInfo, db *sql.DB) error {
	log.Printf("User <%s %s> is new. Adding user to database...", userInfo.GivenName, userInfo.FamilyName)

	_, err := db.Exec(`
        INSERT INTO users (userId, email, givenName, familyName)
        VALUES (?, ?, ?, ?)`,
		userInfo.Id, userInfo.Email, userInfo.GivenName, userInfo.FamilyName,
	)
	if err != nil {
		return fmt.Errorf("Failed adding new user to database\nError:\n%v\n\n", err)
	}

	log.Printf("Successfully added new user <%s %s> to database.", userInfo.GivenName, userInfo.FamilyName)
	return nil
}

func IsNewUser(db *sql.DB, userInfo *UserInfo) error {
	var (
		userId      string
		email       string
		givenName   string
		familyName  string
		joined      time.Time
		poopTotal   int
		failedTotal int
	)

	err := db.QueryRow(`SELECT * FROM users WHERE userId = ?`, userInfo.Id).Scan(
		&userId, &email, &givenName, &familyName, &joined, &poopTotal, &failedTotal,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return addUser(userInfo, db)
		} else {
			return err
		}
	}
	log.Printf("User <%s %s> already exists.", userInfo.GivenName, userInfo.FamilyName)
	return nil
}
