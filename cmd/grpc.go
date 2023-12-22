package main

import (
	"context"
	"errors"
	"golinkcut/api/proto"
	"golinkcut/internal/config"
	"golinkcut/internal/link"
	"google.golang.org/grpc"
	"log"
	"net"
	"os/signal"
	"syscall"
)

func runGrpcServer(uc link.UseCase, cfg config.Config) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	srv := proto.NewGrpcServer(uc)
	lis, err := net.Listen("tcp", ":"+cfg["grpcPort"].(string))
	if err != nil {
		log.Printf("CRITICAL: GRPC server could not listen to port %s!")
		panic(err)
	}
	s := grpc.NewServer()
	proto.RegisterLinkServiceServer(s, &srv)
	log.Printf("GRPC server started at port %s", cfg["grpcPort"].(string))
	go func() {
		if err := s.Serve(lis); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			log.Printf("GRPC Server error: %s", err)
		}
	}()
	// Listen for the interrupt signal.
	<-ctx.Done()
	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	log.Println("shutting down GRPC server gracefully, press Ctrl+C again to force")
	s.GracefulStop()
	log.Println("GRPC server exiting")
}
