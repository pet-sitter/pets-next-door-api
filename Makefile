## Variables ##
BUILD_DIR=bin
SERVER_BINARY_NAME=server
MIGRATE_BINARY_NAME=migrate

SERVER_PORT=8080
SERVER_URL=http://localhost:${SERVER_PORT}
SWAGGER_ROUTE=/swagger/index.html

DATABASE_USER=postgres
DATABASE_PASSWORD=postgres
DATABASE_HOST=localhost
DATABASE_PORT=5454
DATABASE_NAME=pets_next_door_api_dev
DATABASE_URL=postgresql://${DATABASE_USER}:${DATABASE_PASSWORD}@${DATABASE_HOST}:${DATABASE_PORT}/${DATABASE_NAME}?sslmode=disable

## Dependencies ##
deps:
	go mod tidy

## Code ##
format:
	go fmt ./...

lint\:install:
	# golangci-lint
	# https://golangci-lint.run/welcome/install/#local-installation
	brew install golangci-lint
lint:
	golangci-lint run

## Version ##
version:
	. ./scripts/version.sh

## Build ##
docs:
	. ./scripts/swagger-gen.sh
docs\:open:
	open ${SERVER_URL}${SWAGGER_ROUTE}
docs\:clean:
	rm -rf pkg/docs

clean:
	go clean
	rm -rf ${BUILD_DIR}
	make docs:clean

build:
	make docs
	go build -o ${BUILD_DIR}/${SERVER_BINARY_NAME} ./cmd/server

run:
	go run ./cmd/server

test:
	make db:test:up
	go test ./... -count=1 -p=1
	make db:test:down

## Database ##
db\:up:
	docker compose -p pets-next-door-api-dev up -d --remove-orphans
db\:down:
	docker compose -p pets-next-door-api-dev down
db\:destroy:
	docker compose -p pets-next-door-api-dev down -v

db\:test\:up:
	docker compose -f docker-compose-test.yml -p pets-next-door-api-test up -d --remove-orphans
db\:test\:down:
	docker compose -f docker-compose-test.yml -p pets-next-door-api-test down
db\:test\:destroy:
	docker compose -f docker-compose-test.yml -p pets-next-door-api-test down -v

## Migrate ##
migrate\:install:
	# golang-migrate CLI
	# https://github.com/golang-migrate/migrate/tree/master/cmd/migrate
	brew install golang-migrate
migrate\:up:
	migrate -path db/migrations -database="${DATABASE_URL}" up
migrate\:down:
	migrate -path db/migrations -database="${DATABASE_URL}" down
migrate\:create:
	migrate create -ext sql -dir db/migrations -seq $(name)
