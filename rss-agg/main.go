package main

import (
	"log"
	"fmt"
	"os"
	"net/http"
	"encoding/json"
	"github.com/joho/godotenv"
)

func h1(w http.ResponseWriter, r *http.Request) {
	resData, err := json.Marshal(struct{
		Name string
		Age int
	}{
		"Dane",
		24,
	})
	if err != nil {
		log.Printf("Error marshaling JSON: %v\n", err)
	}
	
	fmt.Printf("Sending data: %s\n", resData)
	_, err = w.Write(resData)
	if err != nil {
		log.Printf("Error writing response: %s\n", err)
	}
}

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
		"GET /": h1, 
		"POST /feeds": handleCreateFeed,
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

