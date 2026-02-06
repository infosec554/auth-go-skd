include .env
export

# Default setup
BINARY_NAME=auth-service

# Main commands
run:
	go run cmd/main.go

build:
	go build -o bin/${BINARY_NAME} cmd/main.go

test:
	go test -v ./...

clean:
	go clean
	rm -f bin/${BINARY_NAME}

# Docker commands
docker-up:
	docker compose up --build -d

docker-down:
	docker compose down

docker-logs:
	docker compose logs -f

# Database Migrations (golang-migrate)
# Usage: make migrate-create name=init_schema
migrate-create:
	migrate create -ext sql -dir migrations -seq $(name)

migrate-up:
	migrate -path migrations -database "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=${POSTGRES_SSL_MODE}" up

migrate-down:
	migrate -path migrations -database "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=${POSTGRES_SSL_MODE}" down

migrate-force:
	migrate -path migrations -database "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=${POSTGRES_SSL_MODE}" force $(version)
