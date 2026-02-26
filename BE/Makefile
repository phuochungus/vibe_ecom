APP_LIST=api-gateway user-service product-service order-service payment-service notification-service

.PHONY: test run-api-gateway run-user-service run-product-service run-order-service run-payment-service run-notification-service docker-up docker-down

test:
	go test ./...

run-api-gateway:
	go run ./cmd/api-gateway

run-user-service:
	go run ./cmd/user-service

run-product-service:
	go run ./cmd/product-service

run-order-service:
	go run ./cmd/order-service

run-payment-service:
	go run ./cmd/payment-service

run-notification-service:
	go run ./cmd/notification-service

docker-up:
	docker compose up --build

docker-down:
	docker compose down -v
