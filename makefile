lint:
	golangci-lint run ./...

run:
	go run cmd/server/main.go -configEnv ./.env

migrate:
	go run cmd/migrator/main.go -typeTask "up" -dsn "postgres://POSTGRES_USER:POSTGRES_PASSWORD@localhost:5447/GO_FIBER_BOILERPLATE_DB?sslmode=disable"

seeder:
	go run cmd/seeder/main.go -typeTask "up" -dsn "postgres://POSTGRES_USER:POSTGRES_PASSWORD@localhost:5447/GO_FIBER_BOILERPLATE_DB?sslmode=disable"
