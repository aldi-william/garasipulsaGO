BINARY_NAME=main

build: tidy dep
	CGO_ENABLED=0 GOOS=linux go build -v -o ${BINARY_NAME}

run:
	./${BINARY_NAME}

build_and_run: build run

clean:
	@go clean
	@rm ${BINARY_NAME}

tidy:
	@go mod tidy

dep:
	@go mod download

docker-up: build
	@docker-compose up --build -d
	@docker ps

docker-down:
	@docker-compose down

.PHONY: setup 
setup:
	@go get -u github.com/swaggo/swag/cmd/swag
	@go install github.com/swaggo/swag/cmd/swag

.PHONY: generate-docs
generate-docs:
	@swag init -g controllers/swagger.go
