package main

import (
	"fmt"
	"os"
	"errors"
	"github.com/joho/godotenv"
)

type Env struct {
	vars map[string]string
}

// Registers a new env var for use
func (e *Env) register(inputVars ...string) error {
	// Make sure there are env vars passed in
	if len(inputVars) == 0 {
		return errors.New("No env vars inputted.")
	}

	// Init map
	e.vars = make(map[string]string)

	// Get env vars
	for _, v := range inputVars {
		e.vars[v] = os.Getenv(v)
		if e.vars[v] == "" {
			return fmt.Errorf("No %s defined in .env.", e.vars[v])
		}
	}

	// Logging result
	fmt.Printf("Loaded environment variables...\n")
	for k, v := range e.vars {
		fmt.Printf("%v: %v\n", k, v)
	}
	fmt.Print("\n")

	return nil
}

func (e *Env) get(varName string) (string, error) {
	if e.vars[varName] == "" {
		return "", fmt.Errorf("No env var with the name %s", varName)
	}
	return e.vars[varName], nil
}

func initEnvVars() (*Env, error) {
	initEnv := Env{}

	// Load in .env
	err := godotenv.Load()
	if err != nil {
		return &initEnv, err
	}

	return &initEnv, nil
}
