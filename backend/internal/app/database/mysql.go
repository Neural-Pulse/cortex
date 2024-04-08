package database

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

// MySQL implementa a interface Database para MySQL.
type MySQL struct{}

// SetupDatabase configura e retorna uma conexão com o banco de dados MySQL.
func (m MySQL) SetupDatabase(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// PingDatabase faz um ping no banco de dados MySQL para verificar se a conexão está ativa.
func (m MySQL) PingDatabase(db *sql.DB) error {
	err := db.Ping()
	if err != nil {
		return err
	}
	return nil
}
