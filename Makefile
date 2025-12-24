.PHONY: postgres
postgres:
	docker run -d --name postgres-16 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret postgres:16-alpine

.PHONY: createdb
createdb:
	docker exec -it postgres-16 createdb --username=root --owner=root bank

.PHONY: dropdb
dropdb:
	docker exec -it postgres-16 dropdb --force bank

.PHONY: migrateup
migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/bank?sslmode=disable" -verbose up

.PHONY: migratedown 
migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/bank?sslmode=disable" -verbose down

.PHONY: sqlc 
sqlc:
	sqlc generate

.PHONY: test 
test:
	go test -v ./...

.PHONY: server
server:
	go run ./cmd/bank/main.go

.PHONY: mock
mock:
	mockgen -destination=internal/mocks/store.go -package=mocks github.com/BruceCompiler/bank/internal/repository/postgres Store