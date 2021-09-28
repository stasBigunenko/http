package model

import "github.com/google/uuid"

//Post model struct

type Post struct {
	Id      uuid.UUID `json:"Id"`
	Author  string    `json:"Author"`
	Message string    `json:"Message"`
}

type MyResponse struct {
	Msg string `json:"msg"`
}
