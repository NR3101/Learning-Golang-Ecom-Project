.PHONY: help build run dev lint format generate-docs generate-graph migrate-up migrate-down docker-up docker-down

help:
	@echo "Makefile commands:"
	@echo "  build          - Build the application"
	@echo "  run            - Run the application"
	@echo "  dev            - Run the application in development mode"
	@echo "  lint           - Run linters on the codebase"
	@echo "  format         - Format the codebase using gofmt and goimports"
	@echo "  generate-docs  - Generate API documentation using swag"
	@echo "  generate-graph - Generate GraphQL schema and resolvers using gqlgen"
	@echo "  migrate-up     - Apply database migrations"
	@echo "  migrate-down   - Rollback database migrations"
	@echo "  docker-up      - Start the application dependency services using Docker"
	@echo "  docker-down    - Stop the application dependency services using Docker"


build:
	@echo "Building all binaries..."
	@mkdir -p bin
	@for cmd in cmd/*/; do \
		if [ -d "$$cmd" ]; then \
			binary=$$(basename $$cmd); \
			echo "Building $$binary..."; \
			go build -o bin/$$binary ./$$cmd; \
		fi; \
	done

run:
	go run ./cmd/api

dev:
	go run ./cmd/api

lint: format
	golangci-lint run ./...

format:
	@gofmt -s -w .
	@goimports -w .

generate-docs:
	mkdir -p docs
	swag init -g cmd/api/main.go -o docs --parseDependency --parseInternal --parseDepth 3 --exclude .git,docs,docker,db -d ./,./internal/server

generate-graph:
	@go get github.com/99designs/gqlgen@v0.17.78
	@go run github.com/99designs/gqlgen generate

migrate-up:
	migrate -path db/migrations -database "postgres://postgres:password@localhost:5432/ecomdb?sslmode=disable" up
migrate-down:
	migrate -path db/migrations -database "postgres://postgres:password@localhost:5432/ecomdb?sslmode=disable" down

docker-up:
	docker compose -f docker/docker-compose.yml up -d
docker-down:
	docker compose -f docker/docker-compose.yml down