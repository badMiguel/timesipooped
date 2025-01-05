package poop

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
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

// TODO do this
func notSignedIn() {
}

func PoopRoute(db *sql.DB) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
			return
		}

		userId, err := r.Cookie("user_id")
		if err != nil && err != http.ErrNoCookie {
			notSignedIn()
			return
		}
		accessToken, err := r.Cookie("access_token")
		if err != nil && err != http.ErrNoCookie {
			notSignedIn()
			return
		}

		log.Println(userId, accessToken)

	})
}

func AddPoop(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func AddFailedPoop(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}

func SubPoop(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}

func SubFailedPoop(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}
