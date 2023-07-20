package Database

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func OpenDBConnetction(connString string) (*sql.DB, error) {
	db, err := sql.Open("mysql", connString)

	if err != nil {
		return nil, err
	}

	return db, nil
}
