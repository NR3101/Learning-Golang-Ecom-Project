.PHONY: help build run dev lint migrate-up migrate-down docker-up docker-down

help:
	@echo "Makefile commands:"
	@echo "  build          - Build the application"
	@echo "  run            - Run the application"
	@echo "  dev            - Run the application in development mode"
	@echo "  lint           - Run linters on the codebase"
	@echo "  migrate-up     - Apply database migrations"
	@echo "  migrate-down   - Rollback database migrations"
	@echo "  docker-up      - Start the application using Docker"
	@echo "  docker-down    - Stop the application using Docker"

build:
	go build -o bin/app ./cmd/api

run:
	go run ./cmd/api

dev:
	go run ./cmd/api

lint:
	golangci-lint run ./...

migrate-up:
	migrate -path db/migrations -database "postgres://postgres:password@localhost:5432/ecomdb?sslmode=disable" up
migrate-down:
	migrate -path db/migrations -database "postgres://postgres:password@localhost:5432/ecomdb?sslmode=disable" down

docker-up:
	docker compose -f docker/docker-compose.yml up -d
docker-down:
	docker compose -f docker/docker-compose.yml down