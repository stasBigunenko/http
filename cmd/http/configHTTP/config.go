package configHTTP

import (
	"os"
)

//http data
type Config struct {
	Port    string
	Grpc    string
	RMQPath string
	RMQLog  string
	RMQPass string
}

func Set() *Config {
	var config Config

	config.Port = os.Getenv("PORT_HTTP")
	if config.Port == "" {
		config.Port = ":8085"
	}

	config.Grpc = os.Getenv("GRPC")
	if config.Grpc == "" {
		config.Grpc = ":9000"
	}

	config.RMQPath = os.Getenv("RMQ_PATH")
	if config.RMQPath == "" {
		config.RMQPath = "localhost:5672/"
	}

	config.RMQLog = os.Getenv("RMQ_LOG")
	if config.RMQLog == "" {
		config.RMQLog = "guest"
	}

	config.RMQPass = os.Getenv("RMQ_PASS")
	if config.RMQPass == "" {
		config.RMQPass = "guest"
	}

	return &Config{
		Port:    config.Port,
		Grpc:    config.Grpc,
		RMQPath: config.RMQPath,
		RMQLog:  config.RMQLog,
		RMQPass: config.RMQPass,
	}
}
