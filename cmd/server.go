package cmd

import (
	"src/http/cmd/config"
	"src/http/storage"
	"src/http/storage/inMemory"
)

type Server struct {
	Config  *config.Config
	Storage storage.Storage
}

func New() *Server {
	return &Server{
	}
}

func (s *Server) ServerConfig() error {
	s.Config = config.Set()
	return nil
}

func (s *Server) StorageServer() error {
	s.Storage = inMemory.New()
	return nil
}

