package main

import (
	"google.golang.org/grpc"
	"log"
	"net"
	pb "src/http/api/proto"
	"src/http/cmd/grpcServer/configGRPC"
	"src/http/pkg/gRPC"
	"src/http/storage"
	"src/http/storage/inMemory"
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
		store = redisDB.New(config.RedisAddr, config.RedisPsw, config.RedisDB)
	}

	//create GRPC server
	s := grpc.NewServer()
	pb.RegisterPostServiceServer(s, gRPC.NewGRPCStore(store))

	log.Printf("GRPC server started on port: %v\n", config.TcpPort)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
