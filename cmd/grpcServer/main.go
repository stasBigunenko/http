package main

import (
	"fmt"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"log"
	"net"
	pb "src/http/api/proto"
	"src/http/cmd/grpcServer/configGRPC"
	"src/http/pkg/gRPC"
	"src/http/storage"
	"src/http/storage/inMemory"
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

	fmt.Println(config)

	var store storage.Storage

	switch config.DbType {
	case "inmemory":
		store = inMemory.New()
	case "redis":
		store = redisDB.New(config.RedisAddr, config.RedisDB)
	case "postgres":
		store, _ = postgres.NewPDB(config.PostgresHost, config.PostgresPort, config.PostgresUser, config.PostgresPsw, config.PostgresDB, config.PostgresSSL)
	}

	fmt.Printf("-----------------------------------------%v\n", store)

	//create GRPC server
	s := grpc.NewServer()
	pb.RegisterPostServiceServer(s, gRPC.NewGRPCStore(store))

	log.Printf("GRPC server started on port: %v\n", config.TcpPort)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
