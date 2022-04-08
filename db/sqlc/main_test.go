package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/betNevS/tinybank/pkg/config"

	_ "github.com/go-sql-driver/mysql"
)

const (
	dbDriver = "mysql"
	dbSource = "root:secret@tcp(127.0.0.1:3306)/tinybank?charset=utf8mb4&parseTime=True&loc=Local"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	config, err := config.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}
	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatalln(err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
