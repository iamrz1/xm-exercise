# Makefile

dep:
	@echo "Running go mod tidy"
	@go mod tidy -v
	@echo "Running go mod download"
	@go mod download -x

run: dep
	@echo "Starting up the app..."
	@go run main.go

up:
	@echo "Starting up the app with dependencies in docker"
	@docker-compose up

down:
	@echo "Shutting down the app  with dependencies in docker"
	@docker-compose down

cleanup:
	@echo "Removing related containers"
	@docker ps -a | grep "xm-exercise" | awk '{print $1}' | xargs docker rm -f -
	@echo "Removing leftover images"
	@docker images -a | grep "xm-exercise" | awk '{print $3}' | xargs docker rmi -f

build:
	@echo "Building docker image"
	@docker-compose build

doc: dep
	@echo "Running ./scripts/swagger.sh"
	@chmod +x ./scripts/swagger.sh
	@./scripts/swagger.sh

test:
	@go test ./..