# Run Application Locally > Read README.md

# Docker (Start)

docker-build:
	docker build -t go-bookstore:latest .

postgres-run:
	docker compose -f docker-compose-postgres.yml up --detach

docker-run: postgres-run
	docker compose -f docker-compose.yml up --detach

run-server:
	DRIVER_NAME=postgres DATABASE_URL=postgres://user:p@ssw0rd@localhost:5432/go-bookstore-db?sslmode=disable PORT=2565 ACCESS_TOKEN=token go run main.go

# Docker (Stop)

postgres-stop:
	docker compose -f docker-compose-postgres.yml stop

# Remove
docker-down:
	docker compose -f docker-compose.yml down	

# Run Unit Test
test-unit: 
	go clean -testcache && go test -v --tags=unit ./...

# Run Integration Test
test-integration:
	docker compose -f docker-compose.test.yml up --build --abort-on-container-exit --exit-code-from it_tests



