package main

import (
	// "encoding/json"
	"log"
	"net/http"
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

func get(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func post(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func main() {
	http.HandleFunc("/get", get)
	http.HandleFunc("/post", post)

	log.Fatal(http.ListenAndServe(":1234", nil))
}
