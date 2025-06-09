FROM golang:1.23.9 AS go_build
WORKDIR /app

RUN apt-get update && apt-get install -y protobuf-compiler && rm -rf /var/lib/apt/lists/*
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
ENV PATH="$PATH:$(go env GOPATH)/bin"

COPY go.mod go.sum /app/
RUN go mod download 

COPY /go/ /app/go/ 
COPY /api/ /app/api/
RUN go generate ./... 

RUN go build -o grpc_server ./go/cmd/server

FROM alpine:latest AS auth
WORKDIR /app 
COPY --from=go_build /app/grpc_server .
CMD ["./grpc_server", "0.0.0.0:7001", "auth"]

FROM alpine:latest AS admin
WORKDIR /app 
COPY --from=go_build /app/grpc_server .
CMD ["./grpc_server", "0.0.0.0:7002", "admin"]

FROM gradle:8.14.0-jdk17 AS web_build
WORKDIR /app
RUN apt-get update && apt-get install -y protobuf-compiler && rm -rf /var/lib/apt/lists/*

COPY /java/ . 
RUN rm ./src/main/proto 
COPY /api/ /app/src/main/proto

RUN gradle build --no-daemon

FROM openjdk:17-jdk-alpine AS web
WORKDIR /app

COPY --from=web_build /app/build/libs/*.jar app.jar
CMD ["java", "-jar", "app.jar"]