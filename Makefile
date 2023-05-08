DB_URL=postgres://root:secret@localhost:5432/reminder?sslmode=disable
DB_URL_TEST=postgres://root:secret@localhost:5432/test_reminder?sslmode=disable

lint:
	gofmt -w .
	golangci-lint run
	go vet ./...

createdb:
	docker exec -it postgres createdb --username=root --owner=root reminder

create_testdb:
	docker exec -it postgres createdb --username=root --owner=root test_reminder

dropdb:
	docker exec -it postgres dropdb --username=root reminder

drop_test_db:
	docker exec -it postgres dropdb --username=root test_reminder

migrateup:
	migrate -database ${DB_URL} -path db/migrations up

migrateup_test:
	migrate -database ${DB_URL_TEST} -path db/migrations up

migratedown:
	migrate -database ${DB_URL} -path db/migrations down

migratedown_test:
	migrate -database ${DB_URL_TEST} -path db/migrations down

db-run:
	docker run --name postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:15-alpine

exec-db:
	docker exec -it postgres psql -U root reminder

run:
	go run cmd/reminder/main.go

run-worker:
	go run cmd/worker/main.go

compose-up:
	docker-compose -f docker-compose.yml up --build

test:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -func coverage.out

coverage:
	go test ./... -coverprofile=coverage.out

coverage-html:
	@$(MAKE) coverage
	go tool cover -html=coverage.out

int_test:
	docker-compose -f docker-compose.test.yml up --build --abort-on-container-exit
	docker-compose -f docker-compose.test.yml down --volumes

#GENERATE MOCKS
mocks:
	go generate ./...

#GENERATE SWAGGER DOCS
swag_gen:
	swag init -g cmd/reminder/main.go

.PHONY: lint, format, createdb, dropdb, migrateup, migrateup_test, migratedown, db-run, exec-db, run, test, coverage, coverage-html, int_test, mocks