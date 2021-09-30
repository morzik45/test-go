run:
	go run ./cmd/main.go  -configFile ./config/config_local.json

docker-build:
	docker-compose build exam-app

docker-up:
	docker-compose up -d

docker-up-db:
	docker-compose up -d db

docker-logs:
	docker-compose logs -f exam-app 

docker-down:
	docker-compose down --remove-orphans

test:
	go test -v ./...

migrate:
	migrate -path ./migrations -database 'postgres://postgres:qwerty@0.0.0.0:5436/postgres?sslmode=disable' up

migrate-rollback:
	migrate -path ./migrations -database 'postgres://postgres:qwerty@0.0.0.0:5436/postgres?sslmode=disable' down
