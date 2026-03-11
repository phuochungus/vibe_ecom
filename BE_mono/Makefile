.PHONY: test run postgres-up postgres-down

test:
	go test ./...

run:
	go run ./cmd/api

postgres-up:
	docker compose -f docker-compose.postgres.yml up -d

postgres-down:
	docker compose -f docker-compose.postgres.yml down -v
