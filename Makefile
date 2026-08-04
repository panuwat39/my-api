APP_NAME := my-api
MAIN_PATH := ./cmd/api
COMPOSE_FILE := deployments/docker/docker-compose.yml

.PHONY: help run fmt vet test check build clean \
	mongo-up mongo-down mongo-status mongo-logs

help:
	@echo "Available commands:"
	@echo "  make run           Run API with .env"
	@echo "  make fmt           Format Go files"
	@echo "  make vet           Run go vet"
	@echo "  make test          Run all tests"
	@echo "  make check         Run format, vet and tests"
	@echo "  make build         Build application"
	@echo "  make clean         Remove build output"
	@echo "  make mongo-up      Start MongoDB"
	@echo "  make mongo-down    Stop MongoDB"
	@echo "  make mongo-status  Show MongoDB status"
	@echo "  make mongo-logs    Follow MongoDB logs"

run:
	@test -f .env || (echo ".env file not found"; exit 1)
	@set -a; . ./.env; set +a; go run $(MAIN_PATH)

fmt:
	@gofmt -w $$(find . -type f -name '*.go' -not -path './vendor/*')

vet:
	@go vet ./...

test:
	@go test ./...

check: fmt vet test

build:
	@mkdir -p bin
	@go build -o bin/$(APP_NAME) $(MAIN_PATH)

clean:
	@rm -rf bin coverage.out

mongo-up:
	@docker compose -f $(COMPOSE_FILE) up -d

mongo-down:
	@docker compose -f $(COMPOSE_FILE) down

mongo-status:
	@docker compose -f $(COMPOSE_FILE) ps

mongo-logs:
	@docker compose -f $(COMPOSE_FILE) logs -f mongodb