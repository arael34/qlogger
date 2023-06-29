package main

import (
	"errors"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Environment struct {
	DatabaseUrl    string
	DatabaseName   string
	AuthHeader     string
	AllowedOrigins *map[string]bool
}

func parseAllowedOrigins(str string) *map[string]bool {
	if str == "" {
		return nil
	}

	allowedOrigins := make(map[string]bool)

	for _, origin := range strings.Split(str, ",") {
		allowedOrigins[origin] = true
	}

	return &allowedOrigins
}

func ValidateEnvironment() (*Environment, error) {
	// In production, the environment variables
	// will be set in the OS env.
	if os.Getenv("ENV") == "PROD" {
		DatabaseUrl := os.Getenv("DATABASE_URL")
		DatabaseName := os.Getenv("DATABASE_NAME")
		AuthHeader := os.Getenv("AUTH_HEADER")
		AllowedOrigins := parseAllowedOrigins(os.Getenv("ALLOWED_ORIGINS"))

		if DatabaseUrl == "" ||
			DatabaseName == "" ||
			AuthHeader == "" ||
			AllowedOrigins == nil {
			return nil, errors.New("failed to parse production env")
		}

		return &Environment{DatabaseUrl, DatabaseName, AuthHeader, AllowedOrigins}, nil
	}

	// Parse .env
	var env map[string]string
	env, err := godotenv.Read()
	if err != nil {
		return nil, errors.New("failed to read local env")
	}

	DatabaseUrl := env["DATABASE_URL"]
	DatabaseName := env["DATABASE_NAME"]
	AuthHeader := env["AUTH_HEADER"]
	AllowedOrigins := parseAllowedOrigins(env["ALLOWED_ORIGINS"])

	if DatabaseUrl == "" ||
		DatabaseName == "" ||
		AuthHeader == "" ||
		AllowedOrigins == nil {
		return nil, errors.New("failed to parse local env")
	}

	return &Environment{DatabaseUrl, DatabaseName, AuthHeader, AllowedOrigins}, nil
}
