package main

import (
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	pb "src/http/api/proto"
	"src/http/cmd/grpcServer/configGRPC"
	"src/http/pkg/gRPC"
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

	//inMemory storage
	//store := inMemory.New()

	fmt.Println(config.RedisAddr, config.RedisPsw, config.RedisDB)

	store := redisDB.New(config.RedisAddr, config.RedisPsw, config.RedisDB)

	//create GRPC server
	s := grpc.NewServer()
	pb.RegisterPostServiceServer(s, gRPC.NewGRPCStore(store))

	log.Println("GRPC server started...", config.TcpPort)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
