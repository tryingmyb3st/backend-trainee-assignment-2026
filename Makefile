include .env
export

up:
	@docker compose up -d --build

down:
	@docker compose down --remove-orphans

seed:
	@docker compose run --rm db-migrations

test-unit:
	@go test -v ./tests/unit/...

test-coverage-html:
	@go test -v ./tests/unit/... -coverprofile=coverage.out -coverpkg=./internal/...
	@go tool cover -html=coverage.out -o coverage.html

test-coverage-total:
	@go test -v ./tests/unit/... -coverprofile=coverage.out -coverpkg=./internal/...
	@go tool cover -func=coverage.out

test-e2e:
	@docker compose -f docker-compose.test.yaml up -d --build
	@echo "waiting database..."
	@sleep 2
	@go test -v ./tests/e2e/...
	@docker compose -f docker-compose.test.yaml down

swagger:
	@docker compose run --rm swagger