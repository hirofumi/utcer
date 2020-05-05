.PHONY: help tidy lint test up down

help: ## Show this help
	@awk -F: '/: *## /{sub(/.*## /,"",$$2);printf"make %-10s %s\n",$$1,$$2}' Makefile

tidy: ## Run go mod tidy and go generate
	go mod tidy
	go generate ./...

lint: ## Run GolangCI-Lint
	golangci-lint run

test: ## Run tests
	go test -race -coverprofile=coverage.txt -covermode=atomic ./...

up: ## Start containers
	docker-compose up -d --build

down: ## Stop containers
	docker-compose down
