.PHONY: infra-up infra-down

migrate-up:
	go run ./cmd/migrate-up

migrate-down:
	go run ./cmd/migrate-down

infra-up:
	docker compose up -d

infra-down:
	docker compose down -v