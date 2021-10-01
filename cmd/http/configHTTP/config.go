package configHTTP

import (
	"os"
)

//http data
type Config struct {
	Port string
	Grpc string
}

func Set() *Config {
	var config Config
	config.Port = os.Getenv("PORT_HTTP")
	config.Grpc = os.Getenv("GRPC")

	if config.Port == "" {
		config.Port = ":8085"
	}

	if config.Grpc == "" {
		config.Grpc = ":9000"
	}

	return &Config{
		Port: config.Port,
		Grpc: config.Grpc,
	}
}
