package main

import (
	"golinkcut/api/proto"
	"golinkcut/internal/config"
	"golinkcut/internal/link"
	"golinkcut/pkg/log"
	"google.golang.org/grpc"
	"net"
)

func runGrpcServer(uc link.UseCase, logger log.Logger, cfg config.Config) {
	srv := proto.NewGrpcServer(uc)
	lis, err := net.Listen("tcp", ":"+cfg["grpcPort"].(string))
	if err != nil {
		panic(err)
	}
	s := grpc.NewServer()
	proto.RegisterLinkServiceServer(s, &srv)
	logger.Infof("GRPC server started at port %s", cfg["grpcPort"].(string))
	if err := s.Serve(lis); err != nil {
		logger.Errorf("GRPC Server error: %s", err)
	}
}
