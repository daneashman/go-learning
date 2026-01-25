package main

import (
	"net/http"
	"encoding/json"
	"log"
	"fmt"
	"io"
)

type Feed struct{
	Id			int		`json:id`
	CreatedAt	string	`json:created_at`
	UpdatedAt	string	`json:updated_at`
	Name		string	`json:name`
	Url			string	`json:url`
}

func httpRes(w http.ResponseWriter, code int, res []byte) {
	w.WriteHeader(code)
	_, err := w.Write(res)
	if err != nil {
		log.Fatal(err)
	}
}

/*
curl -X POST -v http://localhost:8000/feeds -d '{"feed_name":"lips", "feed_url": "lips.com"}'
*/
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
	VALUES (DEFAULT, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, '%s', '%s')
	RETURNING *;
	`
	q := fmt.Sprintf(sqlInsertFeed, body.FeedName, body.FeedUrl)
	insertedFeedRow, err := db.Query(q)
	if err != nil {
		httpRes(w, 400, []byte(fmt.Sprintf("%v", err)))
		return
	}

	// Get ID from db response and return to user
	if insertedFeedRow.Next() {
		var insertedFeed Feed
		err = insertedFeedRow.Scan(&insertedFeed.Id, &insertedFeed.CreatedAt, &insertedFeed.UpdatedAt, &insertedFeed.Name, &insertedFeed.Url)
		if err != nil {
			httpRes(w, 500, []byte(fmt.Sprintf("%v", err)))
			return
		}
		
		jsonFeed, err := json.Marshal(insertedFeed)
		if err != nil {
			httpRes(w, 500, []byte(fmt.Sprintf("%v", err)))
			return
		}

		httpRes(w, 201, jsonFeed)
	} else {
		httpRes(w, 500, []byte("Error getting row back from db."))
	}
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
	var feed Feed
	feeds := make([]Feed, 0, 12)
	for feedRows.Next() == true {
		err = feedRows.Scan(&feed.Id, &feed.CreatedAt, &feed.UpdatedAt, &feed.Name, &feed.Url)
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
