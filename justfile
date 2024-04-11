default: up

up:
  air go run *.go

lint:
  golangci-lint run
