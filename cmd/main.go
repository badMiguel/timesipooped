package main

import (
	"log"
	"net/http"
	"os"

	"timesipooped.fyi/internal"

	// todo remove on production
	"github.com/joho/godotenv"
)

type PoopList struct {
	p []string
}

type PoopInfo struct {
	poopTotal  int
	fPoopTotal int

	poopList  *PoopList
	fPoopList *PoopList
}

type User struct {
	Username string    `json:"username"`
	Info     *PoopInfo `json:"info"`
}

func newPoopInfo() *PoopInfo {
	return &PoopInfo{
		poopTotal:  0,
		fPoopTotal: 0,
		poopList:   &PoopList{[]string{}},
		fPoopList:  &PoopList{[]string{}},
	}
}

func newUser(u string) *User {
	return &User{
		Username: u,
		Info:     newPoopInfo(),
	}
}

func handleGetReq(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func handlePostReq(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		panic(err)
	}

	authConf := internal.NewOAuthConfig()
	http.HandleFunc("/login", internal.HandleLogin(authConf))
	http.HandleFunc("/login/callback", internal.HandleCallback(authConf))

	http.HandleFunc("/get", handleGetReq)
	http.HandleFunc("/post", handlePostReq)

	log.Fatal(http.ListenAndServe(os.Getenv("SERVER_PORT"), nil))
}
