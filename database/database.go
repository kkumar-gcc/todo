package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"

	"github.com/kkumar-gcc/todo/constants"
)

// GetInstance returns a singleton instance of the database connection.
func GetInstance() *sql.DB {
	db, err := sql.Open("sqlite3", constants.SqliteDatabaseName)
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	return db
}
