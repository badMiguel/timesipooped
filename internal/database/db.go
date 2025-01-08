package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

    _ "modernc.org/sqlite"
)

type UserSchema struct {
	userId      string
	picture     string
	joined      time.Time
	poopTotal   int
	failedTotal int
}

type UserInfo struct {
	Id            string `json:"id"`
	Picture       string `json:"picture"`
	VerifiedEmail bool   `json:"verified_email"`
}

type PoopSchema struct {
	poopId  int
	date    time.Time
	userId  string
	success bool
}

// user info for client
type ClientUserInfo struct {
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
        picture TEXT,
        accessToken TEXT,
        refreshToken TEXT,
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
        userId TEXT,
        success INTEGER,
        FOREIGN KEY(userId) REFERENCES users(userId) ON DELETE CASCADE
    )`)
	if err != nil {
		return fmt.Errorf("Error creating new poop table: %v\n", err)
	}
	log.Println("Initialized poop table")

	return nil
}

func InitializeDB(db *sql.DB) {
	for i := 0; i < 10; i++ {
		if err := initialize(db); err == nil {
			log.Println("Successfully connected to the database")
			return
		} else {
			log.Println(err)
		}
		log.Printf("Attempt %d failed. Trying again...\n", i)
		if i == 9 {
			panic("Maximum attempts reached. Could not connect to DB")
		}
		time.Sleep(time.Duration(i/2) * time.Second)
	}
}

func addUser(db *sql.DB, userInfo *UserInfo, accessToken, refreshToken string) error {
	log.Printf(
		"User <%s> is new. Adding user to database...",
		userInfo.Id,
	)

	for i := 0; i < 10; i++ {
		_, err := db.Exec(`
        INSERT INTO users (userId, picture, accessToken, refreshToken)
        VALUES (?, ?, ?, ?)`,
			userInfo.Id, userInfo.Picture, accessToken, refreshToken,
		)
		if err != nil {
			log.Printf("Error adding new user to DB... Trying again (%d)", i)
			if i == 9 {
				return fmt.Errorf("Failed adding new user to database: %v\n", err)
			}
			time.Sleep(time.Duration(i/2) * time.Second)
			continue
		}
		break
	}

	log.Printf(
		"Successfully added new user <%s> to database.",
		userInfo.Id,
	)
	return nil
}

func updateDetails(db *sql.DB, userInfo *UserInfo, accessToken, refreshToken string) error {
	log.Printf("User <%s> already exists.", userInfo.Id)

	for i := 0; i < 10; i++ {
		_, err := db.Exec(`
            UPDATE users 
            SET picture = ?,
                accessToken = ?,
                refreshToken = ?
            WHERE userId=?`,
			userInfo.Picture, accessToken, refreshToken, userInfo.Id,
		)
		if err != nil {
			log.Printf("Error updating user details... Trying again (%d)", i)
			if i == 9 {
				return fmt.Errorf("Failed updating user data to database: %v\n", err)
			}
			time.Sleep(time.Duration(i/2) * time.Second)
			continue
		}
		break
	}
	log.Printf(
		"Successfully updated user <%s> data",
		userInfo.Id,
	)
	return nil
}

func IsNewUser(db *sql.DB, userInfo *UserInfo, accessToken, refreshToken string) error {
	var userId string
	err := db.QueryRow(
		`SELECT userId FROM users WHERE userId = ?`, userInfo.Id,
	).Scan(&userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return addUser(db, userInfo, accessToken, refreshToken)
		} else {
			return err
		}
	}
	return updateDetails(db, userInfo, accessToken, refreshToken)
}

func FetchData(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
			return
		}

		userIdCookie, err := r.Cookie("user_id")
		if err != nil {
			log.Printf("Error getting user_id cookie: %v\n", err)
			http.Error(w, "Failed to retrieve user id", http.StatusUnauthorized)
			return
		}

		var (
			picture     string
			poopTotal   int
			failedTotal int
		)
		err = db.QueryRow(`
            SELECT picture, poopTotal, failedTotal
            FROM users 
            WHERE userId = ?`,
			userIdCookie.Value,
		).Scan(
			&picture, &poopTotal, &failedTotal,
		)
		if err != nil {
			log.Printf("Error checking user <%s> data: %v\n", userIdCookie.Value, err)
			http.Error(w, "Failed to get user data", http.StatusInternalServerError)
			return
		}

		clientUserInfo := &ClientUserInfo{
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
