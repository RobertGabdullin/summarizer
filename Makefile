MIGRATIONS_DIR=backend/migrations
USER=postgres
PASSWORD=postgres
DB_NAME=summarizer
HOST=localhost

BACKEND_PATH=backend/cmd/main.go

start:
	go run $(BACKEND_PATH) 

migrate-up:
	goose -dir $(MIGRATIONS_DIR) postgres "user=$(USER) password=$(PASSWORD) dbname=$(DB_NAME) host=$(HOST) sslmode=disable" up

