package main

import (
	"net/http"
	"encoding/json"
	"log"
	"os"
	"fmt"
	"io"
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type ReqBody struct {
	FeedName	string	`json:"feed_name"`
	FeedUrl		string	`json:"feed_url"`
}

func httpRes(w http.ResponseWriter, code int, res []byte) {
	w.WriteHeader(code)
	_, err := w.Write(res)
	if err != nil {
		log.Fatal(err)
	}
}

func handleCreateFeed(w http.ResponseWriter, req *http.Request) {
	// Connect to db
	db, err := sql.Open("pgx", os.Getenv("DB_URI"))
	if err != nil {
		httpRes(w, 500, []byte(fmt.Sprintf("%v", err)))
		return
	}
	defer db.Close()

	// Pull data from req body
	bodyRaw, err := io.ReadAll(req.Body)
	if err != nil {
		httpRes(w, 500, []byte(fmt.Sprintf("%v", err)))
		return
	}

	// Unmarshal req body into ReqBody struct
	var body ReqBody
	err = json.Unmarshal(bodyRaw, &body)
	if err != nil {
		httpRes(w, 500, []byte(fmt.Sprintf("%v", err)))
		return
	}

	// Insert data into db
	sqlInsertFeed := `
	INSERT INTO feeds
	VALUES (DEFAULT, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, '%s', '%s');
	`
	q := fmt.Sprintf(sqlInsertFeed, body.FeedName, body.FeedUrl)
	_, err = db.Query(q)
	if err != nil {
		httpRes(w, 400, []byte(fmt.Sprintf("%v", err)))
		return
	}
}
