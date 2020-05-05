SUBMODULES := pgxutc pqutc test

.PHONY: help tidy lint test up down

help: ## Show this help
	@awk -F: '/: *## /{sub(/.*## /,"",$$2);printf"make %-10s %s\n",$$1,$$2}' Makefile

tidy: ## Run go mod tidy and go generate
	go mod tidy
	go generate ./...
	for d in $(SUBMODULES); do cd "$$d"; go mod tidy || exit 1; go generate ./... || exit 1; cd ..; done

lint: ## Run GolangCI-Lint
	golangci-lint run

test: ## Run tests
	go test -race -coverprofile=coverage.txt -covermode=atomic ./...
	for d in $(SUBMODULES); do cd "$$d"; go test ./... || exit 1; cd ..; done

up: ## Start containers
	docker-compose up -d --build

down: ## Stop containers
	docker-compose down
