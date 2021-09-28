package configGRPC

import (
	"os"
)

//http gRPC data
type Config struct {
	//Host    string
	TcpPort string

	//Storage type
	DbType string

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

	config.DbType = os.Getenv("DB_TYPE")
	if config.DbType == "" {
		config.DbType = "redis"
	}

	config.RedisAddr = os.Getenv("REDIS_ADDR")
	if config.RedisAddr == "" {
		config.RedisAddr = "127.0.0.1:6379"
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
		TcpPort:   config.TcpPort,
		DbType:    config.DbType,
		RedisAddr: config.RedisAddr,
		RedisPsw:  config.RedisPsw,
		RedisDB:   config.RedisDB,
	}
}
