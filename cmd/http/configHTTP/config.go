package configHTTP

import (
	"os"
)

//http data
type Config struct {
	//Host string
	Port string
	//Conn string
	Grpc string
}

func Set() *Config {
	var config Config
	//config.Host = os.Getenv("HOST")
	config.Port = os.Getenv("PORT")
	//config.Conn = os.Getenv("CONN")
	config.Grpc = os.Getenv("GRPC")

	//if config.Host == "" {
	//	config.Host = "127.0.0.1"
	//}
	if config.Port == "" {
		config.Port = ":8085"
	}

	//if config.Conn == "" {
	//	config.Conn = "grpc"
	//}

	if config.Grpc == "" {
		config.Grpc = ":9000"
	}

	return &Config{
		//Host: config.Host,
		Port: config.Port,
		//Conn: config.Conn,
		Grpc: config.Grpc,
	}
}
