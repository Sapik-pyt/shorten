package service

import (
	"context"

	"github.com/Sapik-pyt/shorten/internal/shorten"
	gen "github.com/Sapik-pyt/shorten/proto/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// 
func (s *ShortenService) CreateShortLink(ctx context.Context, req *gen.CreateShortLinkRequest) (*gen.CreateShortLinkResponse, error) {
	if len(req.OriginalLink) == 0 {
		return nil, status.Error(codes.InvalidArgument, "empty link in request")
	}

	shortLink, err := shorten.HashString(req.OriginalLink)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "creating short link: %s", err.Error())
	}

	ok, err := s.repository.CheckExistance(ctx, shortLink)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "checking existance: %s", err.Error())
	}
	if ok {
		return nil, status.Error(codes.AlreadyExists, "short link already exists")
	}

	err = s.repository.Save(ctx, shortLink, req.OriginalLink)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "saving short link: %s", err.Error())
	}

	return &gen.CreateShortLinkResponse{
		ShortLink: shortLink,
	}, nil
}

// 
func (s *ShortenService) FetchOriginalLink(ctx context.Context, req *gen.FetchOriginalLinkRequest) (*gen.FetchOriginalLinkResponse, error) {
	if len(req.ShortLink) == 0 {
		return nil, status.Error(codes.InvalidArgument, "empty link in request")
	}

	originalLink, err := s.repository.Get(ctx, req.ShortLink)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "fetching original link: %s", err.Error())
	}
	if originalLink == nil {
		return nil, status.Error(codes.NotFound, "original link not found")
	}
	return &gen.FetchOriginalLinkResponse{
		OriginalLink: *originalLink,
	}, nil
}
