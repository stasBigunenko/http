package main

import (
	"context"
	"github.com/gorilla/mux"
	"log"
	"os"
	"os/signal"
	"src/http/cmd"
	"src/http/pkg/handlers"
	"src/http/pkg/services"
)

//Start server and initialize server's config, storage, router
func main() {
	server := cmd.New()
	err := server.ServerConfig()
	if err != nil {
		log.Fatal(err)
	}

	store, err := server.StorageServer()
	if err != nil {
		log.Fatal(err)
	}

	services := services.NewService(*store)

	postroutes := handlers.NewHandler(&services)

	r := mux.NewRouter().StrictSlash(false)
	sub := r.PathPrefix("/posts").Subrouter()

	router := postroutes.Routes(sub)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		oscall := <-c
		log.Printf("system call:%+v", oscall)
		cancel()
	}()

	if err := server.Run(ctx, router); err != nil {
		log.Printf("failed to serve:+%v\n", err)
	}
}
