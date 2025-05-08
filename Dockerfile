FROM golang:1.23.9
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

RUN go build -o auth_server ./go/cmd/auth
CMD ["./auth_server", "localhost:7676"]