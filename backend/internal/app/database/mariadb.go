package database

import (
	"github.com/jmoiron/sqlx"
	"fmt"
	_ "github.com/go-sql-driver/mysql" 
)

type MariaDB struct {
	DB *sqlx.DB
}

func SetupMariaDB(user, password, host, dbName string) (*MariaDB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", user, password, host, dbName)
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

return &MariaDB{DB: db}, nil
}