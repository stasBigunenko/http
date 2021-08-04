package config

import (
	"os"
)

//config data
type Config struct {
	Host     string
	Port     string
	}

func Set() *Config {
	var config Config
	config.Host = os.Getenv("HOST")
	config.Port = os.Getenv("PORT")

	if config.Host == "" {
		config.Host = "127.0.0.1"
	}
	if config.Port == "" {
		config.Port = ":8080"
	}

	return &Config{
		Host:     config.Host,
		Port:     config.Port,
	}
}
