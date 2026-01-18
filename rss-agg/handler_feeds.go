package main

import (
	"net/http"
	"encoding/json"
	"log"
	"os"
	"fmt"
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func httpRes(w http.ResponseWriter, code int, res []byte) {
	w.WriteHeader(code)
	_, err := w.Write(res)
	if err != nil {
		log.Fatal(err)
	}
}

func handleCreateFeed(w http.ResponseWriter, _ *http.Request) {
	// Connect to db
	db, err := sql.Open("pgx", os.Getenv("DB_URI"))
	if err != nil {
		httpRes(w, 500, []byte(fmt.Sprintf("%v", err)))
		return
	}
	defer db.Close()

	// Insert data into db
	sqlInsertFeed := `
	INSERT INTO feeds
	VALUES (DEFAULT, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, '%s', '%s');
	`
	q := fmt.Sprintf(sqlInsertFeed, "one", "one.com")
	_, err = db.Query(q)
	if err != nil {
		httpRes(w, 400, []byte(fmt.Sprintf("%v", err)))
		return
	}
}
