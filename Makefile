run:
	go run cmd/auth/main.go --config=./config/local.yaml

compose-up:
	docker-compose -p taqwa-auth up -d