package proto

import (
	"context"
	"errors"
	"golinkcut/internal/link"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type grpcServer struct {
	UnimplementedLinkServiceServer
	uc link.UseCase
}

func (s *grpcServer) CreateLink(ctx context.Context, req *CreateLinkRequest) (*CreateLinkResponse, error) {
	url := req.GetUrl()
	l, err := s.uc.Create(ctx, link.CreateLinkRequest{OriginalLink: url})
	if err != nil {
		var errLinkExists *link.ErrLinkExists
		if errors.As(err, &errLinkExists) {
			return nil, status.Errorf(codes.AlreadyExists, "This link already exists: %v", url)
		}
		return nil, status.Errorf(codes.Unknown, "Error: %v", err)
	}
	res := &CreateLinkResponse{Alias: l.Alias}
	return res, nil
}

func (s *grpcServer) GetLink(ctx context.Context, req *GetLinkRequest) (*GetLinkResponse, error) {
	alias := req.GetAlias()
	l, err := s.uc.Get(ctx, alias)
	if err != nil {
		var errNotExists *link.ErrNotExists
		if errors.As(err, &errNotExists) {
			return nil, status.Errorf(codes.NotFound, "Not Found")
		}
		return nil, status.Errorf(codes.Unknown, "Error: %v", err)
	}
	res := &GetLinkResponse{Url: l.Original}
	return res, nil
}

func NewGrpcServer(uc link.UseCase) grpcServer {
	gs := grpcServer{
		uc: uc,
	}
	return gs
}
