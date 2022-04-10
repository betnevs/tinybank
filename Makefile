mysql:
	docker run --name mysql -e MYSQL_ROOT_PASSWORD=secret -p 3306:3306 -d amd64/mysql --character-set-server=utf8mb4 --collation-server=utf8mb4_unicode_ci

createdb:
	docker exec -it mysql mysql -h127.0.0.1 -uroot -psecret -e "CREATE DATABASE tinybank;"

dropdb:
	docker exec -it mysql mysql -h127.0.0.1 -uroot -psecret -e "DROP DATABASE tinybank;"

migrateup:
	 migrate -path db/migration -database "mysql://root:secret@tcp(127.0.0.1:3306)/tinybank" -verbose up

migrateup1:
	 migrate -path db/migration -database "mysql://root:secret@tcp(127.0.0.1:3306)/tinybank" -verbose up 1

migratedown:
	 migrate -path db/migration -database "mysql://root:secret@tcp(127.0.0.1:3306)/tinybank" -verbose down

migratedown1:
	migrate -path db/migration -database "mysql://root:secret@tcp(127.0.0.1:3306)/tinybank" -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/betNevS/tinybank/db/sqlc Store

.PHONY: mysql createdb dropdb migrateup migratedown sqlc server mock