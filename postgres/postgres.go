package postgres

import (
	"database/sql"
	_ "github.com/lib/pq"
)

func connectToPostgres() (*sql.DB, error) {
	connectionString := "user={USER} dbname={DBNAME} sslmode=verify-full"
	var err error
	db, err := sql.Open("postgres", connectionString)

	if err != nil {
		return nil, err
	}

	return db, nil

}