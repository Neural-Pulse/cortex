package database

import (
	"database/sql"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestDatabaseConnection(t *testing.T) {
	dsn := "srv-cortex:12345678@tcp(localhost:3306)/"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	t.Log("Successfully connected to the database")
}
