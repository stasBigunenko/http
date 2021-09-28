package main

import (
	"google.golang.org/grpc"
	"log"
	"net"
	pb "src/http/api/proto"
	"src/http/pkg/gRPC"
	"src/http/storage/inMemory"
)

func main() {
	//Set http
	cfg := SetGRPC()

	//start listening on tcp
	lis, err := net.Listen("tcp", cfg.TcpPort)
	if err != nil {
		log.Fatalf("failed to listen: %s", err)
	}

	store := inMemory.New()

	//create GRPC server
	s := grpc.NewServer()
	pb.RegisterPostServiceServer(s, gRPC.NewGRPCStore(store))

	log.Println("GRPC storage server starting...")

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
