package service

import (
	"context"
	"errors"
	"testing"

	mock "github.com/Sapik-pyt/shorten/internal/mocks"
	gen "github.com/Sapik-pyt/shorten/proto/gen"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestCreateShortLink_Positive(t *testing.T) {
	ctx := context.Background()
	testCases := []struct {
		name    string
		request *gen.CreateShortLinkRequest
		prepare func(repoMock *mock.MockRepository)
		want    *gen.CreateShortLinkResponse
	}{
		{
			name: "Positive Case",
			request: &gen.CreateShortLinkRequest{
				OriginalLink: "https://www.youtube.com",
			},
			prepare: func(repoMock *mock.MockRepository) {
				repoMock.EXPECT().CheckExistance(ctx, gomock.Any()).Return(false, nil)
				repoMock.EXPECT().Save(ctx, gomock.Any(), gomock.Any()).Return(nil)
			},
			want: &gen.CreateShortLinkResponse{
				ShortLink: "PSQgQZOvzT",
			},
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()
			repoMock := mock.NewMockRepository(controller)
			srv := NewShortenService(repoMock)

			tt.prepare(repoMock)
			response, err := srv.CreateShortLink(ctx, tt.request)
			assert.NoError(t, err, "unexpected error")

			assert.Equal(t, tt.want, response, "not equal")
		})
	}
}

func TestCreateShortLink_Negative(t *testing.T) {
	ctx := context.Background()
	testCases := []struct {
		name          string
		request       *gen.CreateShortLinkRequest
		prepare       func(repoMock *mock.MockRepository)
		wantErrorCode codes.Code
	}{
		{
			name: "Empty original link",
			request: &gen.CreateShortLinkRequest{
				OriginalLink: "",
			},
			prepare:       func(repoMock *mock.MockRepository) {},
			wantErrorCode: codes.InvalidArgument,
		},
		{
			name: "Link already exists",
			request: &gen.CreateShortLinkRequest{
				OriginalLink: "https://www.youtube.com",
			},
			prepare: func(repoMock *mock.MockRepository) {
				repoMock.EXPECT().CheckExistance(ctx, gomock.Any()).Return(true, nil)
			},
			wantErrorCode: codes.AlreadyExists,
		},
		{
			name: "Fail to Save",
			request: &gen.CreateShortLinkRequest{
				OriginalLink: "https://www.youtube.com",
			},
			prepare: func(repoMock *mock.MockRepository) {
				repoMock.EXPECT().CheckExistance(ctx, gomock.Any()).Return(false, nil)
				repoMock.EXPECT().Save(ctx, gomock.Any(), gomock.Any()).Return(errors.New("error"))
			},
			wantErrorCode: codes.Internal,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()
			repoMock := mock.NewMockRepository(controller)
			srv := NewShortenService(repoMock)

			tt.prepare(repoMock)
			_, err := srv.CreateShortLink(ctx, tt.request)
			st, _ := status.FromError(err)
			assert.Equal(t, st.Code(), tt.wantErrorCode, "not equal")
		})
	}
}

func TestFetchOriginalLink_Positive(t *testing.T) {
	st := "https://www.youtube.com"
	ctx := context.Background()
	testCases := []struct {
		name    string
		request *gen.FetchOriginalLinkRequest
		prepare func(repoMock *mock.MockRepository)
		want    *gen.FetchOriginalLinkResponse
	}{
		{
			name: "Positive Case",
			request: &gen.FetchOriginalLinkRequest{
				ShortLink: "PSQgQZOvzT",
			},
			prepare: func(repoMock *mock.MockRepository) {
				repoMock.EXPECT().Get(ctx, gomock.Any()).Return(&st, nil)
			},
			want: &gen.FetchOriginalLinkResponse{
				OriginalLink: "https://www.youtube.com",
			},
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()
			repoMock := mock.NewMockRepository(controller)
			srv := NewShortenService(repoMock)

			tt.prepare(repoMock)
			response, err := srv.FetchOriginalLink(ctx, tt.request)
			assert.NoError(t, err, "no fetch original link")

			assert.Equal(t, tt.want, response, "not equal")
		})
	}
}

func TestFetchOriginalLink_Negative(t *testing.T) {
	ctx := context.Background()
	testCases := []struct {
		name          string
		request       *gen.FetchOriginalLinkRequest
		prepare       func(repoMock *mock.MockRepository)
		wantErrorCode codes.Code
	}{
		{
			name: "Empty link",
			request: &gen.FetchOriginalLinkRequest{
				ShortLink: "",
			},
			prepare:       func(repoMock *mock.MockRepository) {},
			wantErrorCode: codes.InvalidArgument,
		},
		{
			name: "Original link not found",
			request: &gen.FetchOriginalLinkRequest{
				ShortLink: "PSQgQZOvzo",
			},
			prepare: func(repoMock *mock.MockRepository) {
				repoMock.EXPECT().Get(ctx, gomock.Any()).Return(nil, nil)
			},
			wantErrorCode: codes.NotFound,
		},
		{
			name: "Getting original link",
			request: &gen.FetchOriginalLinkRequest{
				ShortLink: "PSQgQZOvzT",
			},
			prepare: func(repoMock *mock.MockRepository) {
				repoMock.EXPECT().Get(ctx, gomock.Any()).Return(nil, errors.New("failed getting original link"))
			},
			wantErrorCode: codes.Internal,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()
			repoMock := mock.NewMockRepository(controller)
			srv := NewShortenService(repoMock)

			tt.prepare(repoMock)
			_, err := srv.FetchOriginalLink(ctx, tt.request)
			st, _ := status.FromError(err)
			assert.Equal(t, st.Code(), tt.wantErrorCode, "not equal")

		})
	}
}
