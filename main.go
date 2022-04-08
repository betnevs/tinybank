package main

import (
	"database/sql"
	"log"

	"github.com/betNevS/tinybank/api"
	db "github.com/betNevS/tinybank/db/sqlc"
	_ "github.com/go-sql-driver/mysql"
)

const (
	dbDriver      = "mysql"
	dbSource      = "root:secret@tcp(127.0.0.1:3306)/tinybank?charset=utf8mb4&parseTime=True&loc=Local"
	serverAddress = "0.0.0.0:8080"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
}
