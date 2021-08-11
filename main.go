package main

import (
	"log"
	"net/http"
	"src/http/cmd"
	"src/http/pkg/handlers"
)

//Start server and initialize server's config
func main() {
	server := cmd.New()
	server.ServerConfig()
	server.StorageServer()
	handler := handlers.New(&server.Storage)
	newRouter := handler.NewRouter()

	log.Fatal(http.ListenAndServe(server.Config.Port, newRouter))
}
