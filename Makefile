build:
	@echo "Building binary..."
	@go build -o bin/api

run: build
	@echo "Running app..."
	@./bin/api

test:
	@echo "Running tests..."
	@go test -v ./...