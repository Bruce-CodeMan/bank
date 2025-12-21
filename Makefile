.PHONY: createdb postgres dropdb migrateup migratedown sqlc test

postgres:
	docker run -d --name postgres-16 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret postgres:16-alpine

createdb:
	docker exec -it postgres-16 createdb --username=root --owner=root bank

dropdb:
	docker exec -it postgres-16 dropdb --force bank

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v ./...

server:
	go run main.go