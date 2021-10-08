package main

import (
	"fmt"
	"log"
	"net"
	"src/http/storage/elasticsearch"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"

	pb "src/http/api/proto"
	"src/http/cmd/grpcServer/configGRPC"
	"src/http/pkg/gRPC"
	"src/http/storage"
	"src/http/storage/inMemory"
	mongoDB "src/http/storage/mongo"
	"src/http/storage/postgres"
	redisDB "src/http/storage/redis"
)

func main() {
	//Set http
	config := configGRPC.SetGRPC()

	//start listening on tcp
	lis, err := net.Listen("tcp", config.TcpPort)
	if err != nil {
		log.Fatalf("failed to listen: %s", err)
	}

	var store storage.Storage

	switch config.DbType {
	case "inmemory":
		store = inMemory.New()
	case "redis":
		store = redisDB.New(config.RedisAddr, config.RedisDB)
	case "postgres":
		store, _ = postgres.NewPDB(config.PostgresHost, config.PostgresPort, config.PostgresUser, config.PostgresPsw, config.PostgresDB, config.PostgresSSL)
	case "mongo":
		store = mongoDB.NewMongo(config.MONGO_INITDB_ROOT_USERNAME, config.MONGO_INITDB_ROOT_PASSWORD, config.MONGO_ADDR)
	case "elastic":
		store, _ = elasticsearch.NewElastic(config.ELK_ADDR)
		if err != nil {
			log.Fatalf("failed to connect elastic: %s", err)
		}
	}

	fmt.Printf("----------------------------------storage - %v-------------------------------\n", config.DbType)

	//create GRPC server
	s := grpc.NewServer()
	pb.RegisterPostServiceServer(s, gRPC.NewGRPCStore(store))

	log.Printf("GRPC server started on port: %v\n", config.TcpPort)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
