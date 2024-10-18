package transport

import (
	"context"
	"database/sql"
	"log"
	"time"

	pb "github.com/duddyV/user-service/proto"
	"github.com/redis/go-redis/v9"
)

type UserServer struct {
	DB    *sql.DB
	Redis *redis.Client
	pb.UserServiceServer
}

func (us *UserServer) CreateUser(ctx context.Context, r *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	query := `INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id`
	var id string
	err := us.DB.QueryRowContext(ctx, query, r.Name, r.Email, r.Password).Scan(&id)
	if err != nil {
		log.Printf("User insert fail: %v", err)
		return nil, err
	}

	rdb := us.Redis
	cacheKey := "user: " + id
	userData := r.Name + ":" + r.Email
	err = rdb.Set(ctx, cacheKey, userData, time.Hour).Err()
	if err != nil {
		log.Printf("Failed to cache user data in Redis: %v", err)
	}

	return &pb.CreateUserResponse{
		Id:    id,
		Name:  r.Name,
		Email: r.Email,
	}, nil
}

// func (us *UserServer) GetUser(ctx context.Context, r *pb.GetUserRequest) (*pb.GetUserResponse, error) {

// }

// func (us *UserServer) UpdateUser(ctx context.Context, r *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {

// }

// func (us *UserServer) DeleteUser(ctx context.Context, r *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {

// }
