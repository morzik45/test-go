run:
	go run ./cmd/main.go  -configfile ./config/config_local.json

build:
	docker-compose build exam-app

up:
	docker-compose up -d

down:
	docker-compose down --remove-orphans

test:
	go test -v ./...

migrate:
	migrate -path ./migrations -database 'postgres://postgres:qwerty@0.0.0.0:5436/postgres?sslmode=disable' up

migrate-rollback:
	migrate -path ./migrations -database 'postgres://postgres:qwerty@0.0.0.0:5436/postgres?sslmode=disable' down
