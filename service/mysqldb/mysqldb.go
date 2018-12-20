package mysqldb

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func init() {
	var err error
	connString := os.Getenv("DATA_SOURCE")
	db, err = sql.Open("mysql", connString) // this does not really open a new connection
	if err != nil {
		log.Fatalf("Error on initializing database connection: %s", err.Error())
	}

	db.SetMaxIdleConns(10)

	err = db.Ping() // This DOES open a connection if necessary. This makes sure the database is accessible
	if err != nil {
		log.Fatalf("Error on opening database connection: %s", err.Error())
	}
}

// Exec just runs a command with mysql. it throws away the return value because
// in our use case we don't need it.
func Exec(statement string) error {
	_, err := db.Exec(statement)
	return err
}
