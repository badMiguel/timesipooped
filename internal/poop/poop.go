package poop

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"timesipooped.fyi/internal/auth"
)

type PoopInfo struct {
	poopTotal  int
	fPoopTotal int

	poopList  []string
	fPoopList []string
}

func newPoopInfo() *PoopInfo {
	return &PoopInfo{
		poopTotal:  0,
		fPoopTotal: 0,
		poopList:   []string{},
		fPoopList:  []string{},
	}
}

func PoopRoute(db *sql.DB) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
			return
		}

		resp, userId, err := auth.VerifyToken(r)
		if err != nil || resp.StatusCode != http.StatusOK {
			log.Printf("Failed to verify if client is logged in")
			http.Error(w, "Failed to authenticate", http.StatusForbidden)
			return
		}

		switch r.URL.Path {
		case "/add":
			AddPoop(w, r, db, userId)
		case "/failed/add":
			AddFailedPoop(w, r, db, userId)
		case "/sub":
			SubPoop(w, r, db, userId)
		case "/failed/sub":
			SubFailedPoop(w, r, db, userId)
		default:
			http.NotFound(w, r)
		}
	})
}

func AddPoop(w http.ResponseWriter, r *http.Request, db *sql.DB, userId string) {
	var poopTotal int
	for i := 0; i < 10; i++ {
		err := db.QueryRow(`SELECT poopTotal FROM users WHERE userId = ?`, userId).Scan(&poopTotal)
		if err != nil {
			log.Printf("Error finding user's <%v> poop total. Trying again(%d)\n", userId, i)
			if i == 9 {
				log.Printf("Max attempt. Cannot query user <%v>\n", userId)
				http.Error(w, "Server cannot find your poop data!", http.StatusInternalServerError)
				return
			}
			time.Sleep(time.Duration(i/2) * time.Second)
			continue
		}
		poopTotal++
		break
	}

	for i := 0; i < 10; i++ {
		_, err := db.Exec(`INSERT INTO poop (userId, success) VALUES (?, ?)`, userId, 1)
		if err != nil {
			log.Printf("Error adding user's <%v> poop. Trying again(%d)\n", userId, i)
			if i == 9 {
				log.Printf("Max attempt. Cannot add user's <%v> poop.\n", userId)
				http.Error(w, "Failed to add your poop :(", http.StatusInternalServerError)
				return
			}
			time.Sleep(time.Duration(i/2) * time.Second)
			continue
		}
		break
	}

	for i := 0; i < 10; i++ {
		_, err := db.Exec(`UPDATE users SET poopTotal = ? WHERE userId = ?`, poopTotal, userId)
		if err != nil {
			log.Printf("Error updating user's <%v> poop total. Trying again(%d)\n", userId, i)
			if i == 9 {
				log.Printf("Max attempt. Cannot add user's <%v> poop.\n", userId)
				http.Error(w, "Failed to add your poop :(", http.StatusInternalServerError)
				return
			}
			time.Sleep(time.Duration(i/2) * time.Second)
			continue
		}
		break
	}
}

func AddFailedPoop(w http.ResponseWriter, r *http.Request, db *sql.DB, userId string) {
	var failedTotal int
	for i := 0; i < 10; i++ {
		err := db.QueryRow(`SELECT failedTotal FROM users WHERE userId = ?`, userId).Scan(&failedTotal)
		if err != nil {
			log.Printf("Error finding user's <%v> failed poop total. Trying again(%d)\n", userId, i)
			if i == 9 {
				log.Printf("Max attempt. Cannot query user <%v>\n", userId)
				http.Error(w, "Server cannot find your failed poop data!", http.StatusInternalServerError)
				return
			}
			time.Sleep(time.Duration(i/2) * time.Second)
			continue
		}
		failedTotal++
		break
	}

	for i := 0; i < 10; i++ {
		_, err := db.Exec(`INSERT INTO poop (userId, success) VALUES (?, ?)`, userId, 0)
		if err != nil {
			log.Printf("Error adding user's <%v> poop. Trying again(%d)\n", userId, i)
			if i == 9 {
				log.Printf("Max attempt. Cannot add user's <%v> poop.\n", userId)
				http.Error(w, "Failed to add your failed poop :(", http.StatusInternalServerError)
				return
			}
			time.Sleep(time.Duration(i/2) * time.Second)
			continue
		}
		break
	}

	for i := 0; i < 10; i++ {
		_, err := db.Exec(`UPDATE users SET failedTotal = ? WHERE userId = ?`, failedTotal, userId)
		if err != nil {
			log.Printf("Error updating user's <%v> failed poop total. Trying again(%d)\n", userId, i)
			if i == 9 {
				log.Printf("Max attempt. Cannot add user's <%v> failed poop.\n", userId)
				http.Error(w, "Failed to add your failed poop :(", http.StatusInternalServerError)
				return
			}
			time.Sleep(time.Duration(i/2) * time.Second)
			continue
		}
		break
	}
}

func SubPoop(w http.ResponseWriter, r *http.Request, db *sql.DB, userId string) {
}

func SubFailedPoop(w http.ResponseWriter, r *http.Request, db *sql.DB, userId string) {
}

//
