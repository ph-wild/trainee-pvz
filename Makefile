BINARY=app

.DEFAULT_GOAL := build

build:
	go build -o $(BINARY) ./cmd

run: build
	./$(BINARY)

test:
	go test -v -cover ./...

clean:
	rm -f $(BINARY)

fmt:
	go fmt ./...

tidy:
	go mod tidy

up:
	docker compose up -d

down:
	docker compose down --rmi all --volumes

migrate-up:
	goose -dir ./migrations postgres "postgres://pvz:pvzpassword@localhost:5445/pvz?sslmode=disable" up

migrate-down:
	goose -dir ./migrations postgres "postgres://pvz:pvzpassword@localhost:5445/pvz?sslmode=disable" down

run-all:
	make up
	@echo "Waiting DB..."
	sleep 5
	make migrate-up
	make run

help:
	@echo "Makefile for Go project:"
	@echo "  make build        - Build the binary"
	@echo "  make run          - Run the application"
	@echo "  make test         - Run tests"
	@echo "  make clean        - Remove built files"
	@echo "  make fmt          - Format the code"
	@echo "  make tidy         - Update dependencies"
	@echo "  make up           - Start DB and prometheus containers"
	@echo "  make down         - Delete DB and prometheus containers with volumes"
	@echo "  make migrate-up   - Goose migrations up"
	@echo "  make migrate-down - Goose migrations down"
	@echo "  make run-all      - Run full preparing steps (up, migrate-up, )"
