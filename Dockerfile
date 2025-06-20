FROM golang:1.23.9 AS go_build
WORKDIR /app
ENV PATH="$PATH:$(go env GOPATH)/bin"
COPY go.mod go.sum /app/
RUN go mod download 
COPY /go/ /app/go/ 
RUN go build -o server ./go/cmd/server

FROM debian:bookworm-slim AS auth
WORKDIR /app 
COPY --from=go_build /app/server .
CMD ["./server", "0.0.0.0:7001", "auth"]

FROM debian:bookworm-slim AS admin
WORKDIR /app 
COPY --from=go_build /app/grpc_server .
CMD ["./server", "0.0.0.0:7002", "admin"]

FROM gradle:8.14.0-jdk17 AS web_build
WORKDIR /app
COPY /java/ . 
RUN gradle build --no-daemon

FROM openjdk:17-jdk-alpine AS web
WORKDIR /app
COPY --from=web_build /app/build/libs/*.jar app.jar
CMD ["java", "-jar", "app.jar"]