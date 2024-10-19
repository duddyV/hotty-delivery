package transport

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"time"

	pb "github.com/duddyV/user-service/proto"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
)

type UserServer struct {
	DB       *sql.DB
	Redis    *redis.Client
	RabbitMQ *amqp.Connection
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

	ch, err := us.RabbitMQ.Channel()
	if err != nil {
		log.Printf("Failed to open a RabbitMQ channel: %v", err)
		return nil, err
	}
	defer ch.Close()

	err = ch.Publish(
		"",
		"user-created",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte("User created with ID: " + id),
		})
	if err != nil {
		log.Printf("Failed to publish a message to RabbitMQ: %v", err)
	}

	return &pb.CreateUserResponse{
		Id:    id,
		Name:  r.Name,
		Email: r.Email,
	}, nil
}

func (us *UserServer) GetUser(ctx context.Context, r *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	cacheKey := "user: " + r.Id
	userData, err := us.Redis.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		query := `SELECT name, email FROM users WHERE id = $1`
		var name, email string
		err := us.DB.QueryRowContext(ctx, query, r.Id).Scan(&name, &email)
		if err != nil {
			return nil, err
		}

		userData = name + ":" + email
		us.Redis.Set(ctx, cacheKey, userData, time.Hour)
	}

	splitData := strings.Split(userData, ":")
	return &pb.GetUserResponse{
		Id:    r.Id,
		Name:  splitData[0],
		Email: splitData[1],
	}, nil
}

func (us *UserServer) UpdateUser(ctx context.Context, r *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	query := `UPDATE users SET name = $1, email = $2 WHERE id = $3`
	_, err := us.DB.ExecContext(ctx, query, r.Name, r.Email, r.Id)
	if err != nil {
		return nil, err
	}

	cacheKey := "user: " + r.Id
	userData := r.Name + ":" + r.Email
	us.Redis.Set(ctx, cacheKey, userData, time.Hour)

	ch, err := us.RabbitMQ.Channel()
	if err != nil {
		return nil, err
	}
	defer ch.Close()

	err = ch.Publish(
		"",
		"user-updated",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte("User updated with ID: " + r.Id),
		})
	if err != nil {
		return nil, err
	}

	return &pb.UpdateUserResponse{}, nil
}

func (us *UserServer) DeleteUser(ctx context.Context, r *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	query := `DELETE FROM users WHERE id = $1`
	_, err := us.DB.ExecContext(ctx, query, r.Id)
	if err != nil {
		return nil, err
	}

	cacheKey := "user: " + r.Id
	us.Redis.Del(ctx, cacheKey)

	ch, err := us.RabbitMQ.Channel()
	if err != nil {
		return nil, err
	}
	defer ch.Close()

	err = ch.Publish(
		"",             // exchange
		"user-deleted", // routing key
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte("User deleted with ID: " + r.Id),
		})
	if err != nil {
		return nil, err
	}

	return &pb.DeleteUserResponse{}, nil
}
