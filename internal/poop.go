package internal

import (
	"net/http"

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

func AddPoop(w http.ResponseWriter, r *http.Request) {
}
