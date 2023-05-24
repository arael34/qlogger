package main

import (
	"errors"

	"github.com/joho/godotenv"
)

type Environment struct {
	DatabaseUrl string
	AuthHeader  string
}

func ValidateEnvironment() (*Environment, error) {
	// Parse .env
	var env map[string]string
	env, err := godotenv.Read()
	if err != nil {
		return &Environment{}, errors.New("failed to read env")
	}

	DatabaseUrl := env["DATABASE_URL"]
	AuthHeader := env["AUTH_HEADER"]

	if DatabaseUrl == "" || AuthHeader == "" {
		return &Environment{}, errors.New("failed to parse env")
	}

	// Return valid environment
	return &Environment{DatabaseUrl, AuthHeader}, nil
}
