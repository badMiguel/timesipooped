package poop

import (
	"database/sql"
	"log"
	"net/http"

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
}

func AddFailedPoop(w http.ResponseWriter, r *http.Request, db *sql.DB, userId string) {
}

func SubPoop(w http.ResponseWriter, r *http.Request, db *sql.DB, userId string) {
}

func SubFailedPoop(w http.ResponseWriter, r *http.Request, db *sql.DB, userId string) {
}
