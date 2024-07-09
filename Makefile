build:
	@echo "Building binary..."
	@go build -o bin/api

run: build
	@echo "Running app..."
	@./bin/api

test:
	@echo "Running tests..."
	@go test -v ./...

run_db:
	@echo "Running DB..."
	@docker run --name mongodb -p 27017:27017 -d mongo:latest

test:
	@echo "Running tests..."
	go test -v ./...