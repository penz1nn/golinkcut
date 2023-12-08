package main

import (
	"golinkcut/api/proto"
	"golinkcut/internal/config"
	"golinkcut/internal/link"
	"golinkcut/pkg/log"
	"google.golang.org/grpc"
	"net"
)

const (
	port = ":50051"
)

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()

	c := config.Config{
		"db":    "memory",
		"debug": true,
	}

	repo := link.NewStorage(c)
	logger := log.NewWithConfig(c)
	uc := link.NewUseCase(repo, logger, c)
	srv := proto.NewGrpcServer(uc)
	proto.RegisterLinkServiceServer(s, &srv)
	logger.Infof("Server started at port %s", port)
	if err := s.Serve(lis); err != nil {
		panic(err)
	}
}
