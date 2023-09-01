DATABASE_URL ?= postgres://user:password@host:port/db-name?sslmode=disable
DOCKER_IMAGE_NAME ?= session_manager_image

.PHONY: migrate docker

migrate:
    migrate -path db/migrations -database "$(DATABASE_URL)" up

docker: migrate
    docker build -t $(DOCKER_IMAGE_NAME) .
    docker run -e DATABASE_URL="$(DATABASE_URL)" $(DOCKER_IMAGE_NAME)