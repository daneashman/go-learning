package main

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"fmt"
	"net/http"
	"encoding/json"
)

// Define env var map globally
// All env vars should be accessed through here, e.g env["PORT"]
var env map[string]string

func initEnvVars (varsToLoad []string) {
	// Make sure there are env vars passed in
	if len(varsToLoad) == 0 {
		log.Fatal("No env vars registered.")
	}

	// Init map for global env map
	env = make(map[string]string)

	// Load in .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	// Get env vars
	for _, v := range varsToLoad {
		env[v] = os.Getenv(v)
		if env[v] == "" {
			log.Fatalf("No %s defined in .env.", env[v])
		}
	}
	
	// Logging result
	fmt.Printf("Loaded environment variables...\n")
	for k, v := range env {
		fmt.Printf("%v: %v\n", k, v)
	}
	fmt.Print("\n")
}

func h1(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Request received from %s.\n", r.Header.Get("User-Agent"))
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
	initEnvVars([]string {
		"PORT",
	})

	// Init handler func
	http.HandleFunc("/", h1)

	// Start listening
	fmt.Printf("Listening on port %s...\n\n", env["PORT"])
	err := http.ListenAndServe(":"+env["PORT"], nil)
	log.Fatal(err)
}

