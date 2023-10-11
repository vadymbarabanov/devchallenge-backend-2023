run:
	go run ./cmd/app/app.go

build:
	go build -o ./bin/app ./cmd/app/app.go

test:
	go test -timeout 30s -v ./internal/...

test-api:
	go test -timeout 30s -v ./cmd/app/...

container:
	docker compose up --build

container-tests:
	docker compose -f ./docker-compose.test.yml up --build