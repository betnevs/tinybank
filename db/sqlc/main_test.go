package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

const (
	dbDriver = "mysql"
	dbSource = "root:secret@tcp(127.0.0.1:3306)/tinybank?charset=utf8mb4&parseTime=True&loc=Local"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error
	testDB, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatalln(err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
