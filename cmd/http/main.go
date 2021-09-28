package main

import (
	"context"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"os"
	"os/signal"
	pb "src/http/api/proto"
	"src/http/cmd/http/configHTTP"
	"src/http/pkg/gRPC/grpccontroller"
	"src/http/pkg/handlers"
	"src/http/pkg/services"
)

//Start server and initialize server's http, storage, router
func main() {
	config := configHTTP.Set()

	//storage := inMemory.New()

	conn, err := grpc.Dial(config.Grpc, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect to grpc: %v", err)
	}
	defer conn.Close()

	store := grpccontroller.New(pb.NewPostServiceClient(conn))

	serv := services.NewService(store)
	postroutes := handlers.NewHandler(&serv)

	r := mux.NewRouter().StrictSlash(false)
	sub := r.PathPrefix("/posts").Subrouter()

	router := postroutes.Routes(sub)

	srv := http.Server{
		Addr:    config.Port,
		Handler: router,
	}

	c := make(chan os.Signal, 1)
	defer close(c)
	signal.Notify(c, os.Interrupt)

	go func() {
		<-c
		srv.Shutdown(context.Background())
	}()

	log.Printf("HTTP server started on port: %v\n", config.Port)

	if err := srv.ListenAndServe(); err != nil {
		log.Printf("failed to serve:+%v\n", err)
	}
}
