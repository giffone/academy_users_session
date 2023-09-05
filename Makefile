DATABASE_URL ?= postgres://user:password@host:port/db-name?sslmode=disable
DOCKER_IMAGE_NAME ?= session_manager_image

.PHONY: migrate build run

migrate:
	migrate -path db/migrations -database "$(DATABASE_URL)" up

build:
	docker build -t $(DOCKER_IMAGE_NAME) .
    
run: migrate build
	docker run -e DATABASE_URL="$(DATABASE_URL)" -p 9090:8080 $(DOCKER_IMAGE_NAME)