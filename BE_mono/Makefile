.PHONY: test run seed-demo postgres-up postgres-down

test:
	go test ./...

run:
	go run ./cmd/api

seed-demo:
	go run ./cmd/seed-demo

postgres-up:
	docker compose -f docker-compose.postgres.yml up -d

postgres-down:
	docker compose -f docker-compose.postgres.yml down -v
