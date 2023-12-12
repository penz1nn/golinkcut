package main

import (
	"golinkcut/api/proto"
	"golinkcut/internal/config"
	"golinkcut/internal/link"
	"google.golang.org/grpc"
	"log"
	"net"
)

func runGrpcServer(uc link.UseCase, cfg config.Config) {
	srv := proto.NewGrpcServer(uc)
	lis, err := net.Listen("tcp", ":"+cfg["grpcPort"].(string))
	if err != nil {
		panic(err)
	}
	s := grpc.NewServer()
	proto.RegisterLinkServiceServer(s, &srv)
	log.Printf("GRPC server started at port %s", cfg["grpcPort"].(string))
	if err := s.Serve(lis); err != nil {
		log.Printf("GRPC Server error: %s", err)
	}
}
