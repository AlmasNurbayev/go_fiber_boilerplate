lint:
	golangci-lint run ./...

run:
	go run cmd/server/main.go -configEnv ./.env