.PHONY: create_migrations migrate_up migrate_down test

create_migration:
	migrate create -dir migrations -ext sql $(NAME)

migrate_up:
	migrate -path migrations -database sqlite3://data/movies.db up

migrate_down:
	migrate -path migrations -database sqlite3://data/movies.db down

database: migrations/*.sql queries/*.sql
	sqlc generate

movies-app: database cmd main.go go.mod go.sum
	go build

test:
	go test ./cmd