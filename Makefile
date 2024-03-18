docker_network:
	docker network create $(network_name)

connect_network:
	docker network connect $(network_name) $(container_name)

postgres:
	docker run --name postgres15 -network authsrvapi-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=password -d postgres:15-alpine

createdb:
	docker exec -it postgres15 createdb --username=root --owner=root users_db

dropdb:
	docker exec -it postgres15 dropdb users_db

migrateup:
	migrate -path db/migration -database "postgresql://root:password@localhost:5432/users_db?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:password@localhost:5432/users_db?sslmode=disable" -verbose down

new_migration:
	migrate create -ext sql -dir db/migration -seq $(name)

sqlc:
	sqlc generate

test:
	go test -v -cover -short ./...

server:
	go run main.go

docker_server_dev:
	 docker run --name authsrvapi -p 8080:8080 authsrvapi:latest

docker_server_prod:
	docker run --name authsrvapi -p 8080:8080 -e GIN_MODE=release authsrvapi:latest

docker_bg_start:
	docker start authsrvapi

docker_bg_stop:
	docker stop authsrvapi

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/techschool/simplebank/db/sqlc Store

.PHONY: docker_network connect_network postgres createdb dropdb migrateup migratedown new_migration sqlc test server docker_server_dev docker_server_prod docker_bg_start docker_bg_stop mock  
