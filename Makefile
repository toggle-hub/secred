-include .env

gen-graphql:
	@go run github.com/99designs/gqlgen generate

run:
	@go run ./cmd/main.go

test:
	@go test -v ./...

migratecreate:
	migrate create -ext sql -dir migrations -seq $(name)

migrateup:
	migrate -path migrations/ -database ${DB_URL} -verbose up

migratedown:
	migrate -path migrations/ -database ${DB_URL} -verbose down 1