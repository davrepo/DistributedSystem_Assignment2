# time_service.Dockerfile

# Use the official Golang image from the Docker Hub
FROM golang:1.20

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod .
COPY go.sum .

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy individual source files and directories
COPY client.log .
COPY readme.md .
COPY running_example.png .

COPY client/client.go ./client/
COPY grpc/proto.pb.go ./grpc/
COPY grpc/proto.proto ./grpc/
COPY grpc/proto_grpc.pb.go ./grpc/
COPY server/server.go ./server/

# Build the server and client
RUN go build -o /app/bin/server ./server/server.go
RUN go build -o /app/bin/client ./client/client.go

# Add a script to serve as an entry point
COPY entrypoint.sh /app/
RUN chmod +x /app/entrypoint.sh

# This form of ENTRYPOINT allows you to pass a command-line argument
ENTRYPOINT ["/app/entrypoint.sh"]
