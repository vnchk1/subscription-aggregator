.PHONY: build run test clean migrate swagger docker-up docker-down docker-clean fmt health

build:
	go build -o bin/subaggregator ./cmd/subaggregator

run:
	go run ./cmd/subaggregator

test:
	go test ./...

clean:
	rm -rf bin/

migrate:
	go run ./cmd/subaggregator

swagger:
	swag init -g cmd/subaggregator/main.go

docker-up:
	docker-compose up --build

docker-down:
	docker-compose down

docker-clean:
	docker-compose down --rmi all --volumes --remove-orphans

lint:
	golangci-lint run

fmt:
	go fmt ./...

health:
	curl http://localhost:8080/health