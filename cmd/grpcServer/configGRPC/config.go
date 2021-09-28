package configGRPC

import (
	"os"
)

//http gRPC data
type Config struct {
	//Host    string
	TcpPort string

	//redis
	RedisAddr string
	RedisPsw  string
	RedisDB   string
}

func SetGRPC() *Config {
	var config Config

	//config.Host = os.Getenv("HOST")
	//if config.Host == "" {
	//	config.Host = "127.0.0.1"
	//}

	config.TcpPort = os.Getenv("TCP_PORT")
	if config.TcpPort == "" {
		config.TcpPort = "127.0.0.1:9000"
	}

	config.RedisAddr = os.Getenv("REDIS_ADDR")
	if config.RedisAddr == "" {
		config.RedisAddr = "localhost:6379"
	}

	config.RedisPsw = os.Getenv("REDIS_PSW")
	if config.RedisPsw == "" {
		config.RedisPsw = "qwerty"
	}

	config.RedisDB = os.Getenv("REDIS_DB")
	if config.RedisDB == "" {
		config.RedisDB = "redisDB"
	}

	return &Config{
		//Host:    config.Host,
		TcpPort:   config.TcpPort,
		RedisAddr: config.RedisAddr,
		RedisPsw:  config.RedisPsw,
		RedisDB:   config.RedisDB,
	}
}
