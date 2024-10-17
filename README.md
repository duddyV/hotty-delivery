# Hotty Delivery Service(SAMPLE)

Hotty Delivery Service is a simple and efficient food delivery platform written in Go. The application allows users to browse restaurants, place orders, and track deliveries in real-time.

## Features

- Browse nearby restaurants
- Place and track orders
- Real-time updates with WebSockets
- Robust backend using gRPC
- Scalable microservice architecture
- Kafka or RabbitMQ for order processing (configurable)

## Tech Stack

- **Backend**: Golang, gRPC
- **Frontend**: React.js, WebSockets
- **Database**: PostgreSQL (SQL-based)
- **Messaging**: Kafka or RabbitMQ (for order handling)
- **Containerization**: Docker

## Prerequisites

- Go 1.20+
- Docker
- PostgreSQL
- Kafka or RabbitMQ

## Installation

1. Clone the repository:

    ```bash
    git clone https://github.com/yourusername/hotty-delivery-service.git
    ```

2. Change into the project directory:

    ```bash
    cd hotty-delivery-service
    ```

3. Build and run the services with Docker:

    ```bash
    docker-compose up --build
    ```

4. Access the app at `http://localhost:3000`

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

