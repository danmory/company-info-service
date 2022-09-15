package app

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"time"

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

func Run() {
	appAddress := os.Getenv("APP_ADDRESS")
	lis, err := net.Listen("tcp", appAddress)
	if err != nil {
		panic("failed to listen on: " + appAddress)
	}
	defer lis.Close()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		oscall := <-c
		log.Printf("system call:%+v", oscall)
		cancel()
	}()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := runGRPCServer(ctx, lis); err != nil {
			log.Println("grpc server is down: " + err.Error())
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := runHTTPProxyServer(ctx, appAddress); err != nil {
			log.Println("proxy server is down: " + err.Error())
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := runSwagger(ctx); err != nil {
			log.Println("proxy server is down: " + err.Error())
		}
	}()
	wg.Wait()
}

func runGRPCServer(ctx context.Context, lis net.Listener) (err error) {
	grpcServer := grpc.NewServer()
	rpc.RegisterCompanyInfoSearcherServer(grpcServer, rpc.NewServer())
	go func() {
		if err = grpcServer.Serve(lis); err != nil && err != http.ErrServerClosed {
			log.Fatalln("grpc server error: " + err.Error())
		}
	}()
	log.Println("grpc server is running")
	<-ctx.Done()
	grpcServer.Stop()
	if err == http.ErrServerClosed {
		err = nil
	}
	log.Println("grpc server stopped")
	return
}

func runHTTPProxyServer(ctx context.Context, gRPCAddress string) (err error) {
	mux := runtime.NewServeMux()
	srv := &http.Server{
		Addr:    os.Getenv("PROXY_ADDRESS"),
		Handler: mux,
	}
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	if err = rest.RegisterCompanyInfoSearcherHandlerFromEndpoint(
		ctx,
		mux,
		gRPCAddress,
		opts); err != nil {
		return
	}
	go func() {
		if err = srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalln("proxy server error: " + err.Error())
		}
	}()
	log.Println("proxy server is running")
	<-ctx.Done()
	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err = srv.Shutdown(ctxShutDown); err != nil {
		log.Fatalln("proxy server shutdown failed: " + err.Error())
	}
	if err == http.ErrServerClosed {
		err = nil
	}
	log.Println("proxy server stopped")
	return

}

func runSwagger(ctx context.Context) (err error) {
	swaggerSpec, err := loads.Embedded(restapi.SwaggerJSON, restapi.FlatSwaggerJSON)
	if err != nil {
		log.Fatalln("failed to run swagger: " + err.Error())
	}
	api := operations.NewCompanyInfoProtoAPI(swaggerSpec)
	server := restapi.NewServer(api)
	server.Port, err = strconv.Atoi(os.Getenv("SWAGGER_PORT"))
	if err != nil {
		log.Fatalln("failed to get port for swagger: " + err.Error())
	}
	go func() {
		if err = server.Serve(); err != nil && err != http.ErrServerClosed {
			log.Fatalln("proxy server error: " + err.Error())
		}
	}()
	log.Println("swagger server is running")
	<-ctx.Done()
	if err = server.Shutdown(); err != nil {
		log.Fatalln("swagger server shutdown failed: " + err.Error())
	}
	if err == http.ErrServerClosed {
		err = nil
	}
	return
}
