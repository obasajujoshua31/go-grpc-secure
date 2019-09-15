package main

import (
	"fmt"
	"go-grpc/api"
	"log"
	"net"
	"strings"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

var (
	port = 7777
)

type contextKey int

const (
	clientIDKey contextKey = iota
)

//

func AuthenticateAgentClient(ctx context.Context, s *api.Server) (string, error) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		clientLogin := strings.Join(md["login"], "")
		clientPassword := strings.Join(md["password"], "")

		if clientLogin != "john" {
			return "", fmt.Errorf("unknown user %s: ", clientLogin)
		}
		if clientPassword != "doe" {
			return "", fmt.Errorf("incorrect password %s :", clientPassword)
		}

		log.Printf("Authenticated Client :%s", clientLogin)
		return "42", nil
	}
	return "", fmt.Errorf("Missing login credentials")
}

func unaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	s, ok := info.Server.(*api.Server)

	if !ok {
		return nil, fmt.Errorf("Unable to cast server")
	}
	clientID, err := AuthenticateAgentClient(ctx, s)

	if err != nil {
		return nil, err
	}
	ctx = context.WithValue(ctx, clientIDKey, clientID)
	return handler(ctx, req)

}

// Application Entry Point

func main() {

	// Create a Listener

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", "localhost", port))
	if err != nil {
		log.Fatal("Failed to listen for connection", err)
	}

	cred, err := credentials.NewServerTLSFromFile("cert/server.crt", "cert/server.key")

	if err != nil {
		log.Fatalf("Could not load TLS keys: %s", err)
	}

	opts := []grpc.ServerOption{grpc.Creds(cred), grpc.UnaryInterceptor(unaryInterceptor)}

	//Create GRPC server
	fmt.Println("Starting GRPC Server ...")
	s := api.Server{}

	// Create a GRPC object
	grpcServer := grpc.NewServer(opts...)

	// Attach the Ping message to the server
	api.RegisterPingServer(grpcServer, &s)

	// Start the server

	fmt.Printf("GRPC Server started at port %s", fmt.Sprintf("%d", port))

	if err = grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to Serve: %s", err)
	}

}
