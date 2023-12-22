package proto

import (
	"context"
	"fmt"
	"golinkcut/internal/config"
	"golinkcut/internal/link"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"net"
	"strings"
	"testing"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func init() {
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	c := config.Config{
		"memory":   true,
		"debug":    false,
		"validate": true,
	}

	repo := link.NewDbStorage(c)
	uc := link.NewUseCase(repo, c)
	srv := NewGrpcService(uc)
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
	url := "yandex.com/?s=golang"
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	client := NewLinkServiceClient(conn)
	resp, err := client.CreateLink(ctx, &CreateLinkRequest{Url: url})
	if err != nil {
		t.Fatalf("CreateLink failed: %v", err)
	}
	gotAlias := resp.GetAlias()
	if len(gotAlias) != 10 {
		t.Fatalf("Wrong alias format: %s", gotAlias)
	}
}

func TestLinkServiceClient_GetLink(t *testing.T) {
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
	gotAlias := resp.GetAlias()
	resp2, err := client.GetLink(context.Background(), &GetLinkRequest{Alias: gotAlias})
	if err != nil {
		t.Fatalf("GetLink failed: %v", err)
	}
	gotUrl := resp2.GetUrl()
	if gotUrl != "google.com" {
		t.Errorf("Got wrong url: %v", gotUrl)
	}
}

func TestLinkServiceClient_CreateLink_Exists(t *testing.T) {
	url := "medium.com/123124/23123018"
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	client := NewLinkServiceClient(conn)
	_, err = client.CreateLink(ctx, &CreateLinkRequest{Url: url})
	if err != nil {
		t.Fatalf("CreateLink failed: %v", err)
	}
	_, err = client.CreateLink(ctx, &CreateLinkRequest{Url: url})
	if err == nil {
		t.Error("Expected not nil error but got nil")
	}
	if !strings.Contains(err.Error(), "code = AlreadyExists") {
		t.Errorf("Wrong error (or error format): %v", err)
	}
}

func TestLinkServiceClient_GetLink_NotExists(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	client := NewLinkServiceClient(conn)
	_, err = client.GetLink(context.Background(), &GetLinkRequest{Alias: "0123456789"})
	if err == nil {
		t.Error("Expected not nil error but got nil")
	}
	if !strings.Contains(err.Error(), "code = NotFound") {
		t.Errorf("Wrong error (or error format): %v", err)
	}
}

func TestLinkServiceClient_CreateLink_BadUrl(t *testing.T) {
	url := "abracadabra-im-not-url"
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	client := NewLinkServiceClient(conn)
	_, err = client.CreateLink(ctx, &CreateLinkRequest{Url: url})
	if err == nil {
		t.Error("Expected not nil error but got nil")
	}
	if !strings.Contains(err.Error(), "code = InvalidArgument") {
		t.Errorf("Wrong error (or error format): %v", err)
	}
}
