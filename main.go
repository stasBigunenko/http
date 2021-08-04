package main

import (
	"log"
	"net/http"
	"src/http/cmd"
	"src/http/pkg/handlers"
)

func main() {
	server := cmd.New()
	server.ServerConfig()
	server.StorageServer()
	store := server.Storage
	handler := handlers.New(&store)
	newRouter := handler.NewRouter()
	log.Fatal(http.ListenAndServe(server.Config.Port, newRouter))
}
