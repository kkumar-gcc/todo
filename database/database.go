package database

import (
	"database/sql"
	"embed"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"

	"github.com/kkumar-gcc/todo/constants"
)

//go:embed migrations/todo_schema.sql
var todoSchema embed.FS

// GetInstance returns a singleton instance of the database connection.
func GetInstance() *sql.DB {
	dbPath, err := GetDatabasePath()
	if err != nil {
		log.Fatal("Failed to get database path:", err)
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	return db
}

// RunMigration reads the schema from the todo_schema.sql file and applies it.
func RunMigration() error {
	db := GetInstance()
	defer db.Close()

	schema, err := todoSchema.ReadFile("migrations/todo_schema.sql")
	if err != nil {
		return fmt.Errorf("failed to read todo_schema.sql file: %v", err)
	}

	_, err = db.Exec(string(schema))
	if err != nil {
		return fmt.Errorf("failed to apply schema: %v", err)
	}

	return nil
}

// GetDatabasePath returns the path to the SQLite database, storing it in a standard location.
func GetDatabasePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	dbDir := filepath.Join(homeDir, ".config", "todo")

	if err := os.MkdirAll(dbDir, os.ModePerm); err != nil {
		return "", err
	}

	return filepath.Join(dbDir, constants.SqliteDatabaseName), nil
}
