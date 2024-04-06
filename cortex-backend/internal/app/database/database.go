package database

import "database/sql"

// Database é a interface que define os métodos comuns para todas as implementações de banco de dados.
type Database interface {
	SetupDatabase(dsn string) (*sql.DB, error)
	PingDatabase(db *sql.DB) error
}
