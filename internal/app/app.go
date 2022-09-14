package app

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/danmory/company-info-service/internal/transport/pb"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Run() {
	lis, err := net.Listen("tcp", os.Getenv("APP_ADDRESS"))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterCompanyInfoSearcherServer(grpcServer, pb.NewServer())
	log.Println("Starting server...")
	go grpcServer.Serve(lis)
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err = pb.RegisterCompanyInfoSearcherHandlerFromEndpoint(context.Background(), mux, os.Getenv("APP_ADDRESS"), opts)
	if err != nil {
		panic(err)
	}
	http.ListenAndServe(":8081", mux)
}

func init() {
	godotenv.Load()
	if os.Getenv("APP_ADDRESS") == "" {
		log.Println("server address did not specified! running on :8080")
		os.Setenv("APP_ADDRESS", "localhost:8080")
	}
}