postgres:
	docker run --name postgres12 --network simple-bank-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=Just2017 -d postgres:14.2-alpine

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres12 dropdb  simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://root:Just2017@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:Just2017@localhost:5432/simple_bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go simple_bank/db/sqlc Store

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test server mock