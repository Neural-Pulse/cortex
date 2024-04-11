// database_config.go

package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

type DatabaseConfig struct {
    ID     int    `db:"id"`
    DSN    string `db:"DSN"`
    DBType string `db:"DBType"`
}

func EnsureDatabaseAndTableExist(db *sqlx.DB, dbName string) error {
	// Ensure the database exists
	createDBQuery := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", dbName)
	if _, err := db.Exec(createDBQuery); err != nil {
		return err
	}

	// Ensure the table exists in the chosen database
	useDBQuery := fmt.Sprintf("USE %s", dbName)
	if _, err := db.Exec(useDBQuery); err != nil {
		return err
	}
	
	createTableQuery := `
		CREATE TABLE IF NOT EXISTS databases_config (
			id INT AUTO_INCREMENT PRIMARY KEY,
			DSN VARCHAR(255) NOT NULL,
			DBType VARCHAR(50) NOT NULL
		) ENGINE=InnoDB`
	if _, err := db.Exec(createTableQuery); err != nil {
		return err
	}
	
	return nil
}

func SaveDatabaseConfig(db *sqlx.DB, config *DatabaseConfig) error {
	_, err := db.NamedExec(`INSERT INTO databases_config (DSN, DBType) VALUES (:DSN, :DBType)`, config)
	return err
}