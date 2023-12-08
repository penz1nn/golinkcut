package proto

import (
	"context"
	"fmt"
	"golinkcut/internal/config"
	"golinkcut/internal/link"
	"golinkcut/pkg/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	defaultlog "log"
	"net"
	"testing"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func init() {
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	c := config.Config{
		"db":    "memory",
		"debug": true,
	}

	repo := link.NewStorage(c)
	logger := log.NewWithConfig(c)
	uc := link.NewUseCase(repo, logger, c)
	srv := NewGrpcServer(uc)
	RegisterLinkServiceServer(s, &srv)
	go func() {
		if err := s.Serve(lis); err != nil {
			panic(fmt.Sprintf("Server exited with error: %v", err))
		}
	}()
}

func bufDialer(ctx context.Context, address string) (net.Conn, error) {
	return lis.Dial()
}

func TestLinkServiceClient_CreateLink(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	client := NewLinkServiceClient(conn)
	resp, err := client.CreateLink(ctx, &CreateLinkRequest{Url: "google.com"})
	if err != nil {
		t.Fatalf("CreateLink failed: %v", err)
	}
	defaultlog.Printf("Got alias: %v", resp.GetAlias())
}
