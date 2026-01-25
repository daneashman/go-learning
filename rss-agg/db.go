package main

import (
	"os"
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func dbConnect() (*sql.DB, error) {
	db, err := sql.Open("pgx", os.Getenv("DB_URI"))
	if err != nil {
		return db, err
	}
	return db, nil
}

