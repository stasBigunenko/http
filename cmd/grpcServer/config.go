package main

import "os"

//http gRPC data
type Config struct {
	Host    string
	TcpPort string
}

func SetGRPC() *Config {
	var config Config
	config.Host = os.Getenv("HOST")
	config.TcpPort = os.Getenv("PORT")

	if config.Host == "" {
		config.Host = "127.0.0.1"
	}
	if config.TcpPort == "" {
		config.TcpPort = ":9000"
	}

	return &Config{
		Host:    config.Host,
		TcpPort: config.TcpPort,
	}
}
