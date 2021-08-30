package cmd

import (
	"context"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"src/http/cmd/config"
	"src/http/storage"
	"src/http/storage/inMemory"
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
		config: config.Set(),
	}
}

func (s *Server) ServerConfig() error {
	s.config = config.Set()
	return nil
}

func (s *Server) StorageServer() (*storage.Storage, error) {
	store := inMemory.New()
	s.storage = store
	return &s.storage, nil
}

func (s *Server) Run(ctx context.Context, router *mux.Router) (err error) {

	log.Println("Server is running on " + s.config.Port)

	srv := &http.Server{
		Addr:    s.config.Port,
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen:%+s\n", err)
		}
	}()

	log.Printf("server started")

	<-ctx.Done()

	log.Printf("server stopped")

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err := srv.Shutdown(ctxShutDown); err != nil {
		log.Fatalf("server Shutdown Failed:%+s", err)
	}

	log.Printf("server exited properly")

	if err == http.ErrServerClosed {
		err = nil
	}

	return
}
