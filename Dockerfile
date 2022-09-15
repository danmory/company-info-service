FROM golang:1.19-alpine AS build

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o go_server github.com/danmory/company-info-service/cmd/app

EXPOSE 8080
EXPOSE 8081
EXPOSE 8082

ENTRYPOINT ["./go_server"]
