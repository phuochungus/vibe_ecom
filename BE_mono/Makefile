.PHONY: test run seed-demo infra-up infra-down postgres-up minio-up

test:
	go test ./...

run:
	go run ./cmd/api

seed-demo:
	go run ./cmd/seed-demo

infra-up:
	docker compose up -d

infra-down:
	docker compose down -v

postgres-up:
	docker compose up -d postgres

minio-up:
	docker compose up -d minio
