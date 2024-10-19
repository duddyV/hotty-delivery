# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY ./services/user-service/go.mod ./services/user-service/go.sum ./
RUN go mod download

COPY ./services/user-service ./

RUN go build -o user-service ./cmd/app

# Runtime stage
FROM alpine:3.18

WORKDIR /app

COPY --from=builder /app/user-service .

CMD ["./user-service"]
