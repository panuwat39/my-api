APP_NAME := my-api
MAIN_PATH := ./cmd/api
COMPOSE_FILE := deployments/docker/docker-compose.yml
ENV_FILE := .env

.PHONY: help run fmt vet test check build clean \
	mongo-up mongo-down mongo-status mongo-logs \
	docker-build docker-up docker-down docker-status docker-logs \
	integration-test

help:
	@echo "Available commands:"
	@echo "  make run              Run API locally"
	@echo "  make fmt              Format Go code"
	@echo "  make vet              Run go vet"
	@echo "  make test             Run tests"
	@echo "  make check            Run fmt, vet, and test"
	@echo "  make build            Build local binary"
	@echo "  make integration-test Run integration tests"
	@echo "  make mongo-up         Start MongoDB"
	@echo "  make mongo-down       Stop MongoDB"
	@echo "  make mongo-status     Show MongoDB status"
	@echo "  make mongo-logs       Follow MongoDB logs"
	@echo "  make docker-build     Build Docker images"
	@echo "  make docker-up        Start API + MongoDB"
	@echo "  make docker-down      Stop API + MongoDB"
	@echo "  make docker-status    Show container status"
	@echo "  make docker-logs      Follow API logs"

run:
	@test -f $(ENV_FILE) || (echo ".env file not found"; exit 1)
	@set -a; . ./$(ENV_FILE); set +a; go run $(MAIN_PATH)

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
	@docker compose --env-file $(ENV_FILE) -f $(COMPOSE_FILE) up -d mongodb

mongo-down:
	@docker compose --env-file $(ENV_FILE) -f $(COMPOSE_FILE) stop mongodb

mongo-status:
	@docker compose --env-file $(ENV_FILE) -f $(COMPOSE_FILE) ps mongodb

mongo-logs:
	@docker compose --env-file $(ENV_FILE) -f $(COMPOSE_FILE) logs -f mongodb

docker-build:
	@docker compose --env-file $(ENV_FILE) -f $(COMPOSE_FILE) build

docker-up:
	@docker compose --env-file $(ENV_FILE) -f $(COMPOSE_FILE) up -d

docker-down:
	@docker compose --env-file $(ENV_FILE) -f $(COMPOSE_FILE) down

docker-status:
	@docker compose --env-file $(ENV_FILE) -f $(COMPOSE_FILE) ps

docker-logs:
	@docker compose --env-file $(ENV_FILE) -f $(COMPOSE_FILE) logs -f api

integration-test:
	@test -f .env.test || (echo ".env.test file not found"; exit 1)
	@set -a; . ./.env.test; set +a; go test -tags=integration -count=1 ./...