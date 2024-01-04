run:
	@go run main.go

build:
	@golangci-lint run
	@go build -o bin/acc.exe .

install:
	@go install .

# Default target
default: build