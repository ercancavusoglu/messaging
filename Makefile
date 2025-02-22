#!make

start:
	@echo Starting project with local mode
	go run cmd/app/main.go

hot:
	@echo Starting project with local air mode
	go install github.com/cosmtrek/air@v1.49.0
	air

dev:
	@echo Starting dev docker compose
	DOCKER_BUILDKIT=0 docker-compose -f docker-compose.yml up -d --build

up:
	docker-compose up

upd:
	docker-compose up -d

logs:
	docker-compose logs -f messaging

# Modules support
tidy:
	go mod tidy

# ==============================================================================
# Linters https://golangci-lint.run/usage/install/

run-linter:
	@echo Starting linters
	golangci-lint run ./...