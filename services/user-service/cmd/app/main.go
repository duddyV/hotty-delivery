package main

import (
	"log"
	"net"
	"os"

	"github.com/duddyV/user-service/internal/connections"
	"github.com/duddyV/user-service/internal/transport"
	pb "github.com/duddyV/user-service/proto"
	"google.golang.org/grpc"
)

func main() {
	// Postgres connection
	db, err := connections.InitPostgres()
	if err != nil {
		log.Fatalf("Postgres connection failed: %v", err)
	}
	defer db.Close()

	// Redis connection
	redisClient, err := connections.InitRedis()
	if err != nil {
		log.Fatalf("Redis connection failed: %v", err)
	}
	defer redisClient.Close()

	// gRPC setup
	port := os.Getenv("USER_SERVICE_PORT")
	if port == "" {
		port = "50051"
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	log.Printf("Listening on port: %s\n", port)

	grpcServer := grpc.NewServer()
	userServer := &transport.UserServer{
		DB:    db,
		Redis: redisClient,
	}
	pb.RegisterUserServiceServer(grpcServer, userServer)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
