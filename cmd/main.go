package main

import (
	"context"
	"net"
	"net/http"
	"os"

	"github.com/Sapik-pyt/shorten/internal/db"
	"github.com/Sapik-pyt/shorten/internal/logging"
	"github.com/Sapik-pyt/shorten/internal/repositories"
	"github.com/Sapik-pyt/shorten/internal/service"
	gen "github.com/Sapik-pyt/shorten/proto/gen"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	grpcPort = ":8088"
	httpPort = ":8080"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Создаем singleton логгера
	logging.CreateLogger()

	adr := runGrpc(ctx)
	runRest(ctx, adr)
}

func runGrpc(ctx context.Context) string {
	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		logging.Logger.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	storage := createStorage(ctx)

	service := service.NewShortenService(storage)
	gen.RegisterShortenServer(s, service)

	go func() {
		logging.Logger.Infof("server listening at %s", grpcPort)

		if err := s.Serve(lis); err != nil {
			logging.Logger.Fatal(err)
		}
	}()

	return lis.Addr().String()
}

func runRest(ctx context.Context, endpointAdr string) {
	mux := runtime.NewServeMux()
	err := gen.RegisterShortenHandlerFromEndpoint(ctx, mux, endpointAdr, []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())})
	if err != nil {
		logging.Logger.Fatal(err)
	}

	logging.Logger.Infof("server listening at %s", httpPort)

	if err := http.ListenAndServe(httpPort, mux); err != nil {
		logging.Logger.Fatal(err)
	}
}

func createStorage(ctx context.Context) service.Repository {
	if os.Getenv("IN_MEMORY") == "true" {
		return repositories.NewInMemoryStorage()
	}
	pool, err := db.ConnectToDb(ctx)
	if err != nil {
		logging.Logger.Fatalf("connecting to db: %s", err)
	}
	return repositories.NewDbStorage(pool)
}
