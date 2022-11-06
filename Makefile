.PHONY: create_migrations migrate_up migrate_down test

create_migration:
	migrate create -dir migrations -ext sql $(NAME)

migrate_up:
	cp data/movies.db data/movies.backup.db
	migrate -path migrations -database sqlite3://data/movies.db up

migrate_down:
	cp data/movies.db data/movies.backup.db
	migrate -path migrations -database sqlite3://data/movies.db down

database: migrations/*.sql queries/*.sql
	sqlc generate

movies-app: database cmd/*.go main.go go.mod go.sum
	go build

test:
	go test ./cmd

letterboxd_export:
	sqlite3 data/movies.db -readonly \
		-init scripts/letterboxd_csv_export.sql