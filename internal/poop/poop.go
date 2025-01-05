package poop

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"timesipooped.fyi/internal/auth"
)

type JsonResponse struct {
	PoopTotal   int `json:"poopTotal"`
	FailedTotal int `json:"failedTotal"`
}

func generateJsonBytes(pTotal int, fTotal int) (*[]byte, error) {
	resp := &JsonResponse{
		PoopTotal:   pTotal,
		FailedTotal: fTotal,
	}
	jsonBytes, err := json.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("Failed to encode json data: %v\n", err)
	}
	return &jsonBytes, nil
}

func generateResponse(w http.ResponseWriter, pTotal, fTotal int) {
	jsonBytes, err := generateJsonBytes(pTotal, fTotal)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed processing data", http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(*jsonBytes)
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

		var poopTotal int
		for i := 0; i < 10; i++ {
			err := db.QueryRow(`SELECT poopTotal FROM users WHERE userId = ?`, userId).Scan(&poopTotal)
			if err != nil {
				log.Printf("Error finding user's <%v> poop total. Trying again(%d)\n", userId, i)
				if i == 9 {
					log.Printf("Max attempt. Cannot query user <%v>: %v\n", userId, err)
					http.Error(w, "Server cannot find your poop data!", http.StatusInternalServerError)
					return
				}
				time.Sleep(time.Duration(i/2) * time.Second)
				continue
			}
			break
		}

		var failedTotal int
		for i := 0; i < 10; i++ {
			err := db.QueryRow(`SELECT failedTotal FROM users WHERE userId = ?`, userId).Scan(&failedTotal)
			if err != nil {
				log.Printf("Error finding user's <%v> failed poop total. Trying again(%d)\n", userId, i)
				if i == 9 {
					log.Printf("Max attempt. Cannot query user <%v>: %v\n", userId, err)
					http.Error(w, "Server cannot find your failed poop data!", http.StatusInternalServerError)
					return
				}
				time.Sleep(time.Duration(i/2) * time.Second)
				continue
			}
			break
		}

		switch r.URL.Path {
		case "/add":
			AddPoop(w, r, db, userId, poopTotal, failedTotal)
		case "/failed/add":
			AddFailedPoop(w, r, db, userId, poopTotal, failedTotal)
		case "/sub":
			SubPoop(w, r, db, userId, poopTotal, failedTotal)
		case "/failed/sub":
			SubFailedPoop(w, r, db, userId, poopTotal, failedTotal)
		default:
			http.NotFound(w, r)
		}
	})
}

func AddPoop(w http.ResponseWriter, r *http.Request, db *sql.DB, userId string, poopTotal, failedTotal int) {
	for i := 0; i < 10; i++ {
		_, err := db.Exec(`INSERT INTO poop (userId, success) VALUES (?, ?)`, userId, 1)
		if err != nil {
			log.Printf("Error adding user's <%v> poop. Trying again(%d)\n", userId, i)
			if i == 9 {
				log.Printf("Max attempt. Cannot add user's <%v> poop: %v\n", userId, err)
				http.Error(w, "Failed to add your poop :(", http.StatusInternalServerError)
				return
			}
			time.Sleep(time.Duration(i/2) * time.Second)
			continue
		}
		break
	}

	for i := 0; i < 10; i++ {
		_, err := db.Exec(`UPDATE users SET poopTotal = ? WHERE userId = ?`, poopTotal+1, userId)
		if err != nil {
			log.Printf("Error updating user's <%v> poop total. Trying again(%d)\n", userId, i)
			if i == 9 {
				log.Printf("Max attempt. Cannot add user's <%v> poop: %v\n", userId, err)
				http.Error(w, "Failed to add your poop :(", http.StatusInternalServerError)
				return
			}
			time.Sleep(time.Duration(i/2) * time.Second)
			continue
		}
		break
	}

	generateResponse(w, poopTotal+1, failedTotal)
}

func AddFailedPoop(w http.ResponseWriter, r *http.Request, db *sql.DB, userId string, poopTotal, failedTotal int) {
	for i := 0; i < 10; i++ {
		_, err := db.Exec(`INSERT INTO poop (userId, success) VALUES (?, ?)`, userId, 0)
		if err != nil {
			log.Printf("Error adding user's <%v> poop. Trying again(%d)\n", userId, i)
			if i == 9 {
				log.Printf("Max attempt. Cannot add user's <%v> poop: %v\n", userId, err)
				http.Error(w, "Failed to add your failed poop :(", http.StatusInternalServerError)
				return
			}
			time.Sleep(time.Duration(i/2) * time.Second)
			continue
		}
		break
	}

	for i := 0; i < 10; i++ {
		_, err := db.Exec(`UPDATE users SET failedTotal = ? WHERE userId = ?`, failedTotal+1, userId)
		if err != nil {
			log.Printf("Error updating user's <%v> failed poop total. Trying again(%d)\n", userId, i)
			if i == 9 {
				log.Printf("Max attempt. Cannot add user's <%v> failed poop: %v\n", userId, err)
				http.Error(w, "Failed to add your failed poop :(", http.StatusInternalServerError)
				return
			}
			time.Sleep(time.Duration(i/2) * time.Second)
			continue
		}
		break
	}

	generateResponse(w, poopTotal, failedTotal+1)
	log.Println(w.Header())
}

func SubPoop(w http.ResponseWriter, r *http.Request, db *sql.DB, userId string, poopTotal, failedTotal int) {
	if poopTotal < 1 {
		http.Error(w, "Stop  - You haven't pooped yet!", http.StatusTeapot)
		return
	}

	for i := 0; i < 10; i++ {
		_, err := db.Exec(`
            DELETE FROM poop WHERE poopId = (
                SELECT poopId
                FROM poop
                WHERE userId = ? AND success = ?
                ORDER BY poopId DESC
                LIMIT 1
            );`, userId, 1,
		)
		if err != nil {
			log.Printf("Error subing user's <%v> poop. Trying again(%d)\n", userId, i)
			if i == 9 {
				log.Printf("Max attempt. Cannot sub user's <%v> poop: %v\n", userId, err)
				http.Error(w, "Failed to subtract your poop :(", http.StatusInternalServerError)
				return
			}
			time.Sleep(time.Duration(i/2) * time.Second)
			continue
		}
		break
	}

	for i := 0; i < 10; i++ {
		_, err := db.Exec(`UPDATE users SET poopTotal = ? WHERE userId = ?`, poopTotal-1, userId)
		if err != nil {
			log.Printf("Error updating user's <%v> poop total. Trying again(%d)\n", userId, i)
			if i == 9 {
				log.Printf("Max attempt. Cannot sub user's <%v> poop: %v\n", userId, err)
				http.Error(w, "Failed to subtract your poop :(", http.StatusInternalServerError)
				return
			}
			time.Sleep(time.Duration(i/2) * time.Second)
			continue
		}
		break
	}

	generateResponse(w, poopTotal-1, failedTotal)
}

func SubFailedPoop(w http.ResponseWriter, r *http.Request, db *sql.DB, userId string, poopTotal, failedTotal int) {
	if failedTotal < 1 {
		http.Error(w, "Stop  - You haven't pooped yet!", http.StatusTeapot)
		return
	}

	for i := 0; i < 10; i++ {
		_, err := db.Exec(`
            DELETE FROM poop WHERE poopId = (
                SELECT poopId
                FROM poop
                WHERE userId = ? AND success = ?
                ORDER BY poopId DESC
                LIMIT 1
            );`, userId, 0,
		)
		if err != nil {
			log.Printf("Error subing user's <%v> failed poop. Trying again(%d)\n", userId, i)
			if i == 9 {
				log.Printf("Max attempt. Cannot sub user's <%v> failed poop: %v\n", userId, err)
				http.Error(w, "Failed to subtract your failed poop :(", http.StatusInternalServerError)
				return
			}
			time.Sleep(time.Duration(i/2) * time.Second)
			continue
		}
		break
	}

	for i := 0; i < 10; i++ {
		_, err := db.Exec(`UPDATE users SET failedTotal = ? WHERE userId = ?`, failedTotal-1, userId)
		if err != nil {
			log.Printf("Error updating user's <%v> failed poop total. Trying again(%d)\n", userId, i)
			if i == 9 {
				log.Printf("Max attempt. Cannot sub user's <%v> failed poop: %v\n", userId, err)
				http.Error(w, "Failed to subtract your failed poop :(", http.StatusInternalServerError)
				return
			}
			time.Sleep(time.Duration(i/2) * time.Second)
			continue
		}
		break
	}

	generateResponse(w, poopTotal, failedTotal-1)
}
