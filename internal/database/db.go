package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type UserSchema struct {
	userId      string
	email       string
	givenName   string
	familyName  string
	picture     string
	joined      time.Time
	poopTotal   int
	failedTotal int
}

type UserInfo struct {
	Id            string `json:"id"`
	Email         string `json:"email"`
	FamilyName    string `json:"family_name"`
	GivenName     string `json:"given_name"`
	Picture       string `json:"picture"`
	VerifiedEmail bool   `json:"verified_email"`
}

type PoopSchema struct {
	poopId int
	date   time.Time
	userId string
}

// user info for client
type ClientUserInfo struct {
	FamilyName  string `json:"family_name"`
	GivenName   string `json:"given_name"`
	Picture     string `json:"picture"`
	PoopTotal   int    `json:"poopTotal"`
	FailedTotal int    `json:"failedTotal"`
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
        picture TEXT,
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

func addUser(db *sql.DB, userInfo *UserInfo) error {
	log.Printf(
		"User <%s> is new. Adding user to database...",
		userInfo.Id,
	)

	_, err := db.Exec(`
        INSERT INTO users (userId, email, givenName, familyName, picture)
        VALUES (?, ?, ?, ?, ?)`,
		userInfo.Id, userInfo.Email, userInfo.GivenName, userInfo.FamilyName, userInfo.Picture,
	)
	if err != nil {
		return fmt.Errorf("Failed adding new user to database\nError:\n%v\n\n", err)
	}

	log.Printf(
		"Successfully added new user <%s> to database.",
		userInfo.Id,
	)
	return nil
}

func updateDetails(db *sql.DB, userInfo *UserInfo) error {
	log.Printf("User <%s> already exists.", userInfo.Id)

	_, err := db.Exec(`
        UPDATE  users SET email=?, givenName=?, familyName=?, picture=? WHERE userId=?`,
		userInfo.Email, userInfo.GivenName, userInfo.FamilyName, userInfo.Picture, userInfo.Id,
	)
	if err != nil {
		return fmt.Errorf("Failed updating user data to database\nError:\n%v\n\n", err)
	}

	log.Printf(
		"Successfully updated user <%s> data",
		userInfo.Id,
	)
	return nil
}

func IsNewUser(db *sql.DB, userInfo *UserInfo) error {
	var (
		userId      string
		email       string
		givenName   string
		familyName  string
		picture     string
		joined      time.Time
		poopTotal   int
		failedTotal int
	)

	err := db.QueryRow(`SELECT * FROM users WHERE userId = ?`, userInfo.Id).Scan(
		&userId, &email, &givenName, &familyName, &picture, &joined, &poopTotal, &failedTotal,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return addUser(db, userInfo)
		} else {
			return err
		}
	}
	return updateDetails(db, userInfo)
}

func FetchData(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIdCookie, err := r.Cookie("user_id")
		if err != nil {
			log.Printf("Error getting user_id cookie \nError:\n%v\n\n", err)
			http.Error(w, "Failed to retrieve user id", http.StatusForbidden)
			return
		}

		var (
			givenName   string
			familyName  string
			picture     string
			poopTotal   int
			failedTotal int
		)
		err = db.QueryRow(`
            SELECT givenName, familyName, picture, poopTotal, failedTotal
            FROM users 
            WHERE userId = ?`,
			userIdCookie.Value,
		).Scan(
			&givenName, &familyName, &picture, &poopTotal, &failedTotal,
		)
		if err != nil {
			log.Printf("Error checking user <%s> data: %v\n", userIdCookie.Value, err)
			http.Error(w, "Failed to get user data", http.StatusInternalServerError)
			return
		}

		clientUserInfo := &ClientUserInfo{
			GivenName:   givenName,
			FamilyName:  familyName,
			Picture:     picture,
			PoopTotal:   poopTotal,
			FailedTotal: failedTotal,
		}
		jsonMessage, err := json.Marshal(clientUserInfo)
		if err != nil {
			log.Printf("Failed to encode clientUserInfo of uid <%s>: %v\n", userIdCookie.Value, err)
			http.Error(w, "Failed encoding json data", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonMessage)
	}
}
