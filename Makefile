DOCKER_COMPOSE_FILE = ./docker-compose.yml

# Build microservices
build:
	@echo "Building services..."
	docker-compose -f $(DOCKER_COMPOSE_FILE) build

# Run containers
up:
	@echo "Starting containers..."
	docker-compose -f $(DOCKER_COMPOSE_FILE) up -d

# Stop containers
down:
	@echo "Stopping containers..."
	docker-compose -f $(DOCKER_COMPOSE_FILE) down

# Full container and data clean
clean:
	@echo "Cleaning up..."
	docker-compose -f $(DOCKER_COMPOSE_FILE) down --volumes --remove-orphans

# Data migrations
migrateup:
	@echo "Creating users table..."
	migrate -path services/user-service/scripts/migrations -database "postgres://postgres:postgres@localhost:5432/hotty_delivery?sslmode=disable" -verbose up

migratedown:
	@echo "Creating users table..."
	migrate -path services/user-service/scripts/migrations -database "postgres://postgres:postgres@localhost:5432/hotty_delivery?sslmode=disable" -verbose down

.PHONY: migrateup migratedown