package model

import (
	"time"
)

//Post model struct
type Post struct {
	Id      int       `json:"Id"`
	Message string    `json:"Message"`
	Time    time.Time `json:"Time"`
}
