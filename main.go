package main

import (
	"context"
	"log"
	"os"
	"os/signal"
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

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		oscall := <-c
		log.Printf("system call:%+v", oscall)
		cancel()
	}()

	if err := server.Run(ctx); err != nil {
		log.Printf("failed to serve:+%v\n", err)
	}
}
