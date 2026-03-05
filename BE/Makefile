BINARY_NAME=stock-management
MAIN_PATH=./cmd/api/main.go
SEED_PATH=./cmd/seed/main.go
MIGRATE_DIR=./db/migrations
DB_DSN?=$(DATABASE_DSN)

.PHONY: run build test lint migrate-up migrate-down migrate-status migrate-force seed sqlc mock tidy air

run:
	go run $(MAIN_PATH)

build:
	go build -o bin/$(BINARY_NAME) $(MAIN_PATH)

test:
	go test ./... -v -race -count=1

test-cover:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out

lint:
	golangci-lint run ./...

migrate-up:
	migrate -path $(MIGRATE_DIR) -database "$(DB_DSN)" up

migrate-down:
	migrate -path $(MIGRATE_DIR) -database "$(DB_DSN)" down

migrate-status:
	migrate -path $(MIGRATE_DIR) -database "$(DB_DSN)" version

migrate-force:
	migrate -path $(MIGRATE_DIR) -database "$(DB_DSN)" force $(VERSION)

seed:
	DATABASE_DSN="$(DB_DSN)" go run $(SEED_PATH)

sqlc:
	sqlc generate

mock:
	go generate ./...

tidy:
	go mod tidy

air:
	air
