package configGRPC

import (
	"os"
)

//http gRPC data
type Config struct {
	//grpc addr
	TcpPort string

	//storage type
	DbType string

	//redis
	RedisAddr string
	RedisDB   string

	//postgres
	PostgresHost string
	PostgresPort string
	PostgresUser string
	PostgresPsw  string
	PostgresDB   string
	PostgresSSL  string

	//mongo
	MONGO_INITDB_ROOT_USERNAME string
	MONGO_INITDB_ROOT_PASSWORD string
	MONGO_ADDR                 string

	//elastic
	ELK_ADDR string
}

func SetGRPC() *Config {
	var config Config

	config.TcpPort = os.Getenv("TCP_PORT")
	if config.TcpPort == "" {
		config.TcpPort = "127.0.0.1:9000"
	}

	config.DbType = os.Getenv("DB_TYPE")
	if config.DbType == "" {
		config.DbType = "inmemory"
	}

	config.RedisAddr = os.Getenv("REDIS_ADDR")
	if config.RedisAddr == "" {
		config.RedisAddr = "127.0.0.1:6379"
	}

	config.RedisDB = os.Getenv("REDIS_DB")
	if config.RedisDB == "" {
		config.RedisDB = "redisDB"
	}

	config.PostgresHost = os.Getenv("POSTGRES_HOST")
	if config.PostgresHost == "" {
		config.PostgresHost = "postgres"
	}

	config.PostgresPort = os.Getenv("POSTGRES_PORT")
	if config.PostgresPort == "" {
		config.PostgresPort = "5432"
	}

	config.PostgresUser = os.Getenv("POSTGRES_USER")
	if config.PostgresUser == "" {
		config.PostgresUser = "postgres"
	}

	config.PostgresPsw = os.Getenv("POSTGRES_PASSWORD")
	if config.PostgresPsw == "" {
		config.PostgresPsw = "qwerty"
	}

	config.PostgresDB = os.Getenv("POSTGRES_DATABASE")
	if config.PostgresDB == "" {
		config.PostgresDB = "postgres"
	}

	config.PostgresSSL = os.Getenv("POSTGRES_SSL")
	if config.PostgresSSL == "" {
		config.PostgresSSL = "disable"
	}

	config.MONGO_INITDB_ROOT_USERNAME = os.Getenv("MONGO_INITDB_ROOT_USERNAME")
	if config.MONGO_INITDB_ROOT_USERNAME == "" {
		config.MONGO_INITDB_ROOT_USERNAME = "root"
	}

	config.MONGO_INITDB_ROOT_PASSWORD = os.Getenv("MONGO_INITDB_ROOT_PASSWORD")
	if config.MONGO_INITDB_ROOT_PASSWORD == "" {
		config.MONGO_INITDB_ROOT_PASSWORD = "root"
	}

	config.MONGO_ADDR = os.Getenv("MONGO_ADDR")
	if config.MONGO_ADDR == "" {
		config.MONGO_ADDR = ":27020"
	}

	config.ELK_ADDR = os.Getenv("ELASTIC_ADDR")
	if config.ELK_ADDR == "" {
		config.ELK_ADDR = "http://localhost:9200"
	}

	return &Config{
		TcpPort: config.TcpPort,
		DbType:  config.DbType,

		RedisAddr: config.RedisAddr,
		RedisDB:   config.RedisDB,

		PostgresHost: config.PostgresHost,
		PostgresPort: config.PostgresPort,
		PostgresUser: config.PostgresUser,
		PostgresPsw:  config.PostgresPsw,
		PostgresDB:   config.PostgresDB,
		PostgresSSL:  config.PostgresSSL,

		MONGO_ADDR:                 config.MONGO_ADDR,
		MONGO_INITDB_ROOT_USERNAME: config.MONGO_INITDB_ROOT_USERNAME,
		MONGO_INITDB_ROOT_PASSWORD: config.MONGO_INITDB_ROOT_PASSWORD,

		ELK_ADDR: config.ELK_ADDR,
	}
}
