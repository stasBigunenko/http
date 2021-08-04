package model

import (
	"time"
)

type Post struct {
	Id    		int 		`json:"Id"`
	Message		string		`json:"Message"`
	Time		time.Time	`json:"Time"`
}

