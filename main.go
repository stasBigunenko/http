package main

import (
	"log"
	"src/http/cmd"
)

//Start server and initialize server's config
func main() {
	server := cmd.New()
	err := server.ServerConfig()
	if err != nil {
		log.Fatal(err)
	}

	err = server.StorageServer()
	if err != nil {
		log.Fatal(err)
	}
	server.ConfigRoutes()

	err = server.Run()
	if err != nil {
		log.Fatal(err)
	}
}
