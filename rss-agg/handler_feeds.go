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

// type Res struct {
// 	message string
// }

// func buildRes(message string) (string, error) {
// 	// Create response struct
// 	res := Res{
// 		message: "message"
// 	}
//
// 	// Marshal to JSON
// 	resJson, err := json.Marshal(res)
// 	if err != nil {
// 		return "", err
// 	}
//
// 	return resJson, nil
// }

func httpRes(w http.ResponseWriter, code int, res []byte) {
	w.WriteHeader(500)
	w.Write([]byte("Error connecting to database."))
}

func handleCreateFeed(w http.ResponseWriter, _ *http.Request) {
	// Connect to db
	db, err := sql.Open("pgx", os.Getenv("DB_URI"))
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error connecting to database."))
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
		w.WriteHeader(400)
		w.Write([]byte(fmt.Sprintf("%v", err)))
	}

	val, err := json.Marshal(struct{
		Name string
		Url string
	}{
		"feed",
		"example.com",
	})
	if err != nil {
		log.Printf("Error marshaling json: %s\n", err)
	}

	_, err = w.Write(val)
	if err != nil {
		log.Printf("Error writing response: %s\n", err)
	}
}
