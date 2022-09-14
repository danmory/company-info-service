package app

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/danmory/company-info-service/internal/docs/restapi"
	"github.com/danmory/company-info-service/internal/docs/restapi/operations"
	rest "github.com/danmory/company-info-service/internal/transport/rest"
	"github.com/danmory/company-info-service/internal/transport/rpc"
	"github.com/go-openapi/loads"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Run() {
	appAddress := os.Getenv("APP_ADDRESS")
	lis, err := net.Listen("tcp", appAddress)
	if err != nil {
		panic("failed to listen on: " + appAddress)
	}
	defer lis.Close()
	var wg sync.WaitGroup
	wg.Add(3) // if at least one service is down - all is down
	go func() {
		defer wg.Done()
		if err := runGRPCServer(lis); err != nil {
			log.Println("grpc server is down: " + err.Error())
		}
	}()
	go func() {
		defer wg.Done()
		if err := runHTTPProxyServer(appAddress); err != nil {
			log.Println("proxy server is down: " + err.Error())
		}
	}()
	go func() {
		defer wg.Done()
		if err := runSwagger(); err != nil {
			log.Println("proxy server is down: " + err.Error())
		}
	}()
	log.Println("Servers are up...")
	wg.Wait()
}

func init() {
	godotenv.Load()
	if os.Getenv("APP_ADDRESS") == "" {
		log.Println("server address did not specified! running on :8080")
		os.Setenv("APP_ADDRESS", "localhost:8080")
	}
	if os.Getenv("PROXY_ADDRESS") == "" {
		log.Println("proxy address did not specified! running on :8081")
		os.Setenv("PROXY_ADDRESS", "localhost:8081")
	}
	if os.Getenv("SWAGGER_PORT") == "" {
		log.Println("swagger port did not specified! running on :8082")
		os.Setenv("SWAGGER_PORT", "localhost:8082")
	}
}

func runGRPCServer(lis net.Listener) error {
	grpcServer := grpc.NewServer()
	rpc.RegisterCompanyInfoSearcherServer(grpcServer, rpc.NewServer())
	return grpcServer.Serve(lis)
}

func runHTTPProxyServer(gRPCAddress string) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	if err := rest.RegisterCompanyInfoSearcherHandlerFromEndpoint(
		ctx,
		mux,
		gRPCAddress,
		opts); err != nil {
		return err
	}
	return http.ListenAndServe(os.Getenv("PROXY_ADDRESS"), mux)
}

func runSwagger() error {
	swaggerSpec, err := loads.Embedded(restapi.SwaggerJSON, restapi.FlatSwaggerJSON)
	if err != nil {
		log.Fatalln(err)
	}

	api := operations.NewCompanyInfoProtoAPI(swaggerSpec)
	server := restapi.NewServer(api)
	defer server.Shutdown()
	server.Port, err = strconv.Atoi(os.Getenv("SWAGGER_PORT"))
	if err != nil {
		return err
	}
	return server.Serve()
}
