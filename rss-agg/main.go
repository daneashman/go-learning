package main

import (
	"log"
	"fmt"
	"net/http"
	"encoding/json"
)

func h1(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Request received at %s from %s.\n", r.Pattern, r.Header.Get("User-Agent"))
	w.Header().Set("Content-Type", "application/json")

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

func main() {
	// Setup and import env vars
	env, err := initEnvVars()
	if err != nil {
		log.Fatal(err)
	}
	env.import("PORT")

	// Register handler funcs
	http.HandleFunc("/", h1)

	// Set the port var
	port, err := env.get("PORT")
	if err != nil {
		log.Fatal(err) 
	}

	// Start listening
	fmt.Printf("Listening on port %s...\n\n", port)
	err = http.ListenAndServe(":"+port, nil)
	log.Fatal(err)
}

