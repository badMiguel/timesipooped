package internal

import (
	// "github.com/mattn/go-sqlite3"
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
