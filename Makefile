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