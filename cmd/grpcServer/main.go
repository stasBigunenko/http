package main

import (
	"google.golang.org/grpc"
	"log"
	"net"
	pb "src/http/api/proto"
	"src/http/cmd/grpcServer/configGRPC"
	"src/http/pkg/gRPC"
	"src/http/storage/inMemory"
)

func main() {
	//Set http
	cfg := configGRPC.SetGRPC()

	//start listening on tcp
	lis, err := net.Listen("tcp", cfg.TcpPort)
	if err != nil {
		log.Fatalf("failed to listen: %s", err)
	}

	//inMemory storage
	store := inMemory.New()

	//redisDB storage
	//redis, err := strconv.Atoi(cfg.RedisDB)
	//if err != nil {
	//	return
	//}
	//
	//store := redisDB.New(cfg.RedisAddr, cfg.RedisPsw, redis)

	//create GRPC server
	s := grpc.NewServer()
	pb.RegisterPostServiceServer(s, gRPC.NewGRPCStore(store))

	log.Println("GRPC server started...", cfg.TcpPort)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
