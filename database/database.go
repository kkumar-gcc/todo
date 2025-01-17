package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

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

// RunMigration reads the schema from the schema.sql file and applies it.
func RunMigration() error {
	db := GetInstance()
	defer db.Close()

	schema, err := os.ReadFile("database/schema.sql")
	if err != nil {
		return fmt.Errorf("failed to read schema.sql file: %v", err)
	}

	_, err = db.Exec(string(schema))
	if err != nil {
		return fmt.Errorf("failed to apply schema: %v", err)
	}

	return nil
}
