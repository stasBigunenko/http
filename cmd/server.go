package cmd

import (
	"github.com/gorilla/mux"
	"src/http/cmd/config"
	"src/http/storage"
	"src/http/storage/inMemory"
)

//Server config
type Server struct {
	Config  *config.Config
	Storage storage.Storage
	Router  *mux.Router
}

func New() *Server {
	return &Server{}
}

func (s *Server) ServerConfig() error {
	s.Config = config.Set()
	return nil
}

func (s *Server) StorageServer() error {
	s.Storage = inMemory.New()
	return nil
}
