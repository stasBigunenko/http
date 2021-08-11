package cmd

import (
	"context"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"src/http/cmd/config"
	"src/http/pkg/handlers"
	"src/http/storage"
	"src/http/storage/inMemory"
	"syscall"
	"time"
)

//Server config
type Server struct {
	config  *config.Config
	storage storage.Storage
	router  *mux.Router
}

func New() *Server {
	return &Server{
		router: mux.NewRouter(),
	}
}

func (s *Server) ServerConfig() error {
	s.config = config.Set()
	return nil
}

func (s *Server) StorageServer() error {
	stor := inMemory.New()
	s.storage = stor
	return nil
}

func (s *Server) ConfigRoutes() {
	postroutes := handlers.New(s.router, s.storage)
	postroutes.Routes()
}

func (s *Server) Run(ctx context.Context) error {
	log.Println("Server is running on " + s.config.Port)
	srv := &http.Server{
		Addr:    s.config.Port,
		Handler: s.router,
	}
	go func() {
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-done
	log.Println("Server Stopped")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := srv.Shutdown(ctx)
	if err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Print("Server Exited Properly")
	return nil
}
