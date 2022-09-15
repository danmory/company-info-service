# Company Info Service

## Description

The service returns info about the company given its INN.

The service consists of 3 servers:

1. gRPC Server
2. HTTP Proxy Server(proxies HTTP requests to gRPC)
3. Swagger Server

## Requirements

## Usage

1. Initialize .env file with the following fields

    ```dotenv
        APP_ADDRESS=0.0.0.0:8080
        PROXY_ADDRESS=0.0.0.0:8081
        SWAGGER_PORT=8082
    ```

2. Run Docker

    `` $ docker build -t company_info_service . ``

    `` $ docker run --rm -d -p 8080:8080 -p 8081:8081 -p 8082:8082 --env-file .env --name company_info_service_server company_info_service ``

3. gRPC runs on APP_ADDRESS (<http://localhost:8080>)

   HTTP Proxy runs on PROXY_ADDRESS (path to request - <http://localhost:8081/inn/{enter_inn}>)

   Swagger runs on SWAGGER_PORT (path with docs - <http://localhost:8082/docs>)

## Additional inforamtion

For gRPC and proxy code re-generation please refer to <https://github.com/grpc-ecosystem/grpc-gateway>

## Contacts

Danila Moriakov - d.moriakov@gmail.com
