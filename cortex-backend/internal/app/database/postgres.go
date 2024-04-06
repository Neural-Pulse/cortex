package database

import (
	"database/sql"

	_ "github.com/lib/pq"
)

// PostgreSQL implementa a interface Database para PostgreSQL.
type PostgreSQL struct{}

// SetupDatabase configura e retorna uma conexão com o banco de dados PostgreSQL.
func (p PostgreSQL) SetupDatabase(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// PingDatabase faz um ping no banco de dados PostgreSQL para verificar se a conexão está ativa.
func (p PostgreSQL) PingDatabase(db *sql.DB) error {
	err := db.Ping()
	if err != nil {
		return err
	}
	return nil
}
