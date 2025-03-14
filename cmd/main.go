package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"timesipooped.fyi/internal/auth"
	"timesipooped.fyi/internal/database"
	"timesipooped.fyi/internal/poop"

    _ "modernc.org/sqlite"

	// todo remove on production
	"github.com/joho/godotenv"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Cookie")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		panic(err)
	}

	db, err := sql.Open("sqlite", "../internal/database/db.sqlite")
	if err != nil {
		panic(fmt.Sprintf("Error opening/creating sqlite3 file: %v\n", err))
	}
	defer db.Close()
	database.InitializeDB(db)

	mux := http.NewServeMux()

	authConf := auth.NewOAuthConfig()
	mux.HandleFunc("/auth/login", auth.HandleLogin(authConf))
	mux.HandleFunc("/auth/login/callback", auth.HandleCallback(authConf, db))
	mux.HandleFunc("/auth/status", auth.HandleStatus(authConf, db))
	mux.HandleFunc("/auth/logout", auth.HandleLogout)

	mux.Handle("/poop/add", http.StripPrefix("/poop", poop.PoopRoute(db, authConf)))
	mux.Handle("/poop/failed/add", http.StripPrefix("/poop", poop.PoopRoute(db, authConf)))
	mux.Handle("/poop/sub", http.StripPrefix("/poop", poop.PoopRoute(db, authConf)))
	mux.Handle("/poop/failed/sub", http.StripPrefix("/poop", poop.PoopRoute(db, authConf)))

	mux.HandleFunc("/get/user", database.FetchData(db))

	log.Fatal(http.ListenAndServe(os.Getenv("SERVER_PORT"), corsMiddleware(mux)))
}
