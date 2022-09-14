package app

import (
	"log"
	"net"
	"os"

	"github.com/danmory/company-info-service/internal/transport/pb"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func Run() {
	lis, err := net.Listen("tcp", os.Getenv("APP_ADDRESS"))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterCompanyInfoSearcherServer(grpcServer, pb.NewServer())
	log.Println("Starting server...")
	grpcServer.Serve(lis)
}

func init() {
	godotenv.Load()
	if os.Getenv("APP_ADDRESS") == "" {
		log.Println("server address did not specified! running on :8080")
		os.Setenv("APP_ADDRESS", "localhost:8080")
	}
}