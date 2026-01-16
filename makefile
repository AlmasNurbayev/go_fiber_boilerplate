lint:
	golangci-lint run ./...

# генерация документации, требуется бинарный файл swag
# go install github.com/swaggo/swag/cmd/swag@v1.16.6
swag:
	swag init -g main.go -d cmd/server,internal/httpApp/handlers --parseInternal --parseDependency -o ./docs/swagger/

run: swag
	go run cmd/server/main.go -configEnv ./.env

migrate:
	go run cmd/migrator/main.go -typeTask "up" -dsn "postgres://POSTGRES_USER:POSTGRES_PASSWORD@localhost:5447/GO_FIBER_BOILERPLATE_DB?sslmode=disable"

seeder:
	go run cmd/seeder/main.go -typeTask "up" -dsn "postgres://POSTGRES_USER:POSTGRES_PASSWORD@localhost:5447/GO_FIBER_BOILERPLATE_DB?sslmode=disable"

