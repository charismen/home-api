# Home API Makefile
# Provides commands for both Docker and local development

.PHONY: help build up down stop restart logs ps \
        run install test clean \
        db-setup db-migrate \
        env-setup

# Default target
help:
	@echo "Home API Makefile"
	@echo ""
	@echo "Docker Commands:"
	@echo "  make build       - Build Docker images"
	@echo "  make up          - Start all Docker containers"
	@echo "  make down        - Stop and remove all Docker containers"
	@echo "  make stop        - Stop all Docker containers without removing them"
	@echo "  make restart     - Restart all Docker containers"
	@echo "  make logs        - View logs from all containers"
	@echo "  make ps          - List running containers"
	@echo ""
	@echo "Local Development Commands:"
	@echo "  make run         - Run the application locally"
	@echo "  make install     - Install dependencies"
	@echo "  make test        - Run tests"
	@echo ""
	@echo "Database Commands:"
	@echo "  make db-setup    - Set up database tables"
	@echo "  make db-migrate  - Run database migrations"
	@echo ""
	@echo "Utility Commands:"
	@echo "  make clean       - Clean build artifacts"
	@echo "  make env-setup   - Create .env file from .env.example"

# Docker commands
build:
	docker-compose build

up:
	docker-compose up -d

down:
	docker-compose down

stop:
	docker-compose stop

restart:
	docker-compose restart

logs:
	docker-compose logs -f

ps:
	docker-compose ps

# Local development commands
run:
	go run cmd/api/main.go

install:
	go mod download
	go mod tidy

test:
	go test -v ./...

# Database commands
db-setup:
	@echo "Setting up database tables..."
	@if [ -z "$(shell docker ps -q -f name=mysql)" ]; then \
		echo "MySQL container is not running. Starting it..."; \
		docker-compose up -d mysql; \
		echo "Waiting for MySQL to be ready..."; \
		sleep 10; \
	fi
	@echo "Database setup complete."

db-migrate:
	@echo "No migrations available yet. Add your migration commands here."

# Utility commands
clean:
	rm -rf bin/
	go clean

env-setup:
	@if [ ! -f .env ]; then \
		echo "Creating .env file from .env.example..."; \
		cp .env.example .env; \
		echo ".env file created. Please update with your configuration."; \
	else \
		echo ".env file already exists."; \
	fi