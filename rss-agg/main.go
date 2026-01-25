package main

import (
	"log"
	"fmt"
	"os"
	"net/http"
	"github.com/joho/godotenv"
)

// Create struct that implements http.Handler
// Allows us to add logging and middlewear
type httpHandler struct {
	f func(w http.ResponseWriter, r *http.Request)
}
func (h httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Request received at %s from %s.\n", r.Pattern, r.Header.Get("User-Agent"))
	w.Header().Set("Content-Type", "application/json")
	h.f(w, r)
}

func main() {
	// Set env vars
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	// Register handler funcs on endpoints
	endpoints := map[string]func(w http.ResponseWriter, r *http.Request){
		"POST /feeds":	handleCreateFeed,
		"GET /feeds":	handleGetFeeds,
	}
	for e, f := range endpoints {
		http.Handle(e, httpHandler{f: f})
	}

	// Start listening
	port := os.Getenv("PORT")
	fmt.Printf("Listening on port %s...\n\n", port)
	err = http.ListenAndServe(":"+port, nil)
	log.Fatal(err)
}

