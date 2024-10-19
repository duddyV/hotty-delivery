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

	// RabbitMQ connection
	rabbitMQConn, err := connections.InitRabbitMQ()
    if err != nil {
        log.Fatalf("Failed to connect to RabbitMQ: %v", err)
    }
    defer rabbitMQConn.Close()

	// gRPC setup
	port := os.Getenv("REDIS_PORT")

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	log.Printf("Listening on port: %s\n", port)

	grpcServer := grpc.NewServer()
	userServer := &transport.UserServer{
		DB:    db,
		Redis: redisClient,
		RabbitMQ: rabbitMQConn,
	}
	pb.RegisterUserServiceServer(grpcServer, userServer)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
