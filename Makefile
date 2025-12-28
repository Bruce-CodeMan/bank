DB_DRIVER 	?= postgres
DB_HOST 	?= localhost
DB_PORT		?= 5432
DB_USER		?= root
DB_PASS		?= secret
DB_NAME		?= bank
DB_SSL		?= disable

POSTGRES_IMAGE	?= postgres:16-alpine
POSTGRES_CONTAINER ?= postgres-16
NETWORK ?= bank-network

MIGRATE_PATH ?= db/migration

DB_SOURCE = postgresql://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSL)

.PHONY: postgres
postgres:
	docker run -d \
		--name $(POSTGRES_CONTAINER) \
		-p $(DB_PORT):5432 \
		--network $(NETWORK) \
		-e POSTGRES_USER=$(DB_USER) \
		-e POSTGRES_PASSWORD=$(DB_PASS) \
		$(POSTGRES_IMAGE)

.PHONY: createdb
createdb:
	docker exec -it $(POSTGRES_CONTAINER) createdb \
		--username=$(DB_USER) \
		--owner=$(DB_USER) \
		$(DB_NAME)

.PHONY: dropdb
dropdb:
	docker exec -it $(POSTGRES_CONTAINER) dropdb  \
	--force $(DB_NAME)

.PHONY: migrateup
migrateup:
	migrate -path $(MIGRATE_PATH) \
			-database "$(DB_SOURCE)" \
			-verbose up

.PHONY: migrateup1
migrateup1:
	migrate -path $(MIGRATE_PATH) \
			-database "$(DB_SOURCE)" \
			-verbose up 1

.PHONY: migratedown 
migratedown:
	migrate -path $(MIGRATE_PATH) \
			-database "$(DB_SOURCE)" \
			-verbose down

.PHONY: migratedown1
migratedown1:
	migrate -path $(MIGRATE_PATH) \
			-database "$(DB_SOURCE)" \
			-verbose down 1

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