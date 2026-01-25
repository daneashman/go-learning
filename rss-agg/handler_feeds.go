package main

import (
	"net/http"
	"encoding/json"
	"log"
	"fmt"
	"io"
)

func httpRes(w http.ResponseWriter, code int, res []byte) {
	w.WriteHeader(code)
	_, err := w.Write(res)
	if err != nil {
		log.Fatal(err)
	}
}

func handleCreateFeed(w http.ResponseWriter, req *http.Request) {
	// Connect to db
	db, err := dbConnect()
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
	type ReqBody struct {
		FeedName	string	`json:"feed_name"`
		FeedUrl		string	`json:"feed_url"`
	}
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

	httpRes(w, 201, []byte(""))
}

func handleGetFeeds(w http.ResponseWriter, req *http.Request) {
	// Connect to db
	db, err := dbConnect()
	if err != nil {
		httpRes(w, 500, []byte(fmt.Sprintf("%v", err)))
		return
	}
	defer db.Close()

	// Pull list of feeds from db
	sqlGetFeeds := `
	SELECT * FROM feeds;
	`
	feedRows, err := db.Query(sqlGetFeeds)
	if err != nil {
		httpRes(w, 500, []byte(fmt.Sprintf("%v", err)))
		return
	}

	// Put together array of feeds
	type Feed struct{
		Id			string
		Created_at	string 
		Updated_at	string
		Name		string
		Url			string
	}
	var feed Feed
	feeds := make([]Feed, 0, 12)
	for feedRows.Next() == true {
		err = feedRows.Scan(&feed.Id, &feed.Created_at, &feed.Updated_at, &feed.Name, &feed.Url)
		if err != nil {
			httpRes(w, 500, []byte(fmt.Sprintf("%v", err)))
			return
		}
		feeds = append(feeds, feed)
	}

	// Marshal array of feeds to JSON
	jsonFeeds, err := json.Marshal(feeds)
	if err != nil {
		httpRes(w, 500, []byte(fmt.Sprintf("%v", err)))
		return
	}

	httpRes(w, 200, jsonFeeds)
}
