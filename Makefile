DB_URL=postgres://root:secret@localhost:5432/reminder?sslmode=disable

createdb:
	docker exec -it postgres createdb --username=root --owner=root reminder

dropdb:
	docker exec -it postgres dropdb --username=root reminder

migrateup:
	migrate -database ${DB_URL} -path db/migrations up

migratedown:
	migrate -database ${DB_URL} -path db/migrations down

db-run:
	docker run --name postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:15-alpine

exec-db:
	docker exec -it postgres psql -U root reminder

.PHONY: createdb, dropdb, migrateup, migratedown, db-run, exec-db