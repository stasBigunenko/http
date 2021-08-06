package model

//Post model struct

type Post struct {
	Id      int    `json:"Id"`
	Author  string `json:"Author"`
	Message string `json:"Message"`
	//Time    time.Time `json:"Time"`
}
