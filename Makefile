.PHONY: up down lint test test-cover

up:
	docker compose up -d

down:
	docker compose down

test:
	go test -race ./pkg/... ./templates/... ./services/...

test-cover:
	go test -race -coverprofile=coverage.out ./pkg/...
	go tool cover -func=coverage.out
	@echo "Detailed HTML report: go tool cover -html=coverage.out"

tidy:
	go mod tidy
	cd templates/rest-service && go mod tidy
	cd templates/worker-service && go mod tidy
