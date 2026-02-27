.PHONY: test run mysql-up mysql-down

test:
	go test ./...

run:
	go run ./cmd/api

mysql-up:
	docker compose -f docker-compose.mysql.yml up -d

mysql-down:
	docker compose -f docker-compose.mysql.yml down -v
