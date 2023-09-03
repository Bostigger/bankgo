postgres:
	docker run --name postgresBank -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=password -d postgres

createdb:
	docker exec -it postgresBank createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgresBank dropdb simple_bank

migrateup:
	migrate -path C:/Users/Silan/GolandProjects/bankgo/db/migrations -database "postgresql://root:password@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path C:/Users/Silan/GolandProjects/bankgo/db/migrations -database "postgresql://root:password@localhost:5432/simple_bank?sslmode=disable" -verbose down

PWD := $(CURDIR)
sqlcgenerate:
	docker run --rm -v $(PWD):/src -w /src kjconroy/sqlc generate

.PHONY:createdb
