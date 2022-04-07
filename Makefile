mysql:
	docker run --name mysql -e MYSQL_ROOT_PASSWORD=secret -p 3306:3306 -d amd64/mysql --character-set-server=utf8mb4 --collation-server=utf8mb4_unicode_ci

createdb:
	docker exec -it mysql mysql -h127.0.0.1 -uroot -psecret -e "CREATE DATABASE tinybank;"

dropdb:
	docker exec -it mysql mysql -h127.0.0.1 -uroot -psecret -e "DROP DATABASE tinybank;"

migrateup:
	 migrate -path db/migration -database "mysql://root:secret@tcp(127.0.0.1:3306)/tinybank" -verbose up

migratedown:
	 migrate -path db/migration -database "mysql://root:secret@tcp(127.0.0.1:3306)/tinybank" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

.PHONY: mysql createdb dropdb migrateup migratedown sqlc