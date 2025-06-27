package mysql

import (
	"database/sql"
	"fmt"

	// Import driver for mysql connectivity.
	_ "github.com/go-sql-driver/mysql"
)

func New(username, password, dbName string) (*sql.DB, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(localhost:3306)/%s", username, password, dbName))
	if err != nil {
		return nil, err
	}

	return db, nil
}
