package models

import (
	"time"

	"github.com/tracewayapp/go-lightning/lpg"
)

type Exception struct {
	Id        int
	Archived  bool
	FirstSeen time.Time
	LastSeen  time.Time
	Title     string
}

func init() {
	lpg.Register[Exception]()
}
