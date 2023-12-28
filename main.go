package main

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/zensu/iota-go/internal/protogen"
	"github.com/zensu/iota-go/internal/service/chat"
	"google.golang.org/grpc"
	"net"
)

func main() {
	// Create a new gRPC server
	grpcServer := grpc.NewServer()

	// Create a new connection pool
	var conn []*chat.Connection

	pool := &chat.PoolConnections{
		Connections: conn,
	}

	// Register the pool with the gRPC server
	protogen.RegisterBroadcastServer(grpcServer, pool)

	// Create a TCP listener at port 8080
	listener, err := net.Listen("tcp", "127.0.0.1:8080")

	if err != nil {
		log.Fatal().Msgf("error creating the server %v", err)
	}

	fmt.Println("Server started at port :8080")

	// Start serving requests at port 8080
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal().Msgf("error creating the server %v", err)
	}
}
