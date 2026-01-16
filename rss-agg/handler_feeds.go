package main

import (
	"net/http"
	"encoding/json"
	"log"
)

func handleCreateFeed(w http.ResponseWriter, _ *http.Request) {
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
