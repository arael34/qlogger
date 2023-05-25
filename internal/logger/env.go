package logger

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type Environment struct {
	DatabaseUrl string
	AuthHeader  string
}

func ValidateEnvironment() (*Environment, error) {
	// In production, the environment variables
	// will be set in the OS env.
	if os.Getenv("ENV") == "PROD" {
		DatabaseUrl := os.Getenv("DATABASE_URL")
		AuthHeader := os.Getenv("AUTH_HEADER")

		if DatabaseUrl == "" || AuthHeader == "" {
			return nil, errors.New("failed to parse production env")
		}

		return &Environment{DatabaseUrl, AuthHeader}, nil
	}

	// Parse .env
	var env map[string]string
	env, err := godotenv.Read()
	if err != nil {
		return nil, errors.New("failed to read local env")
	}

	DatabaseUrl := env["DATABASE_URL"]
	AuthHeader := env["AUTH_HEADER"]

	if DatabaseUrl == "" || AuthHeader == "" {
		return nil, errors.New("failed to parse local env")
	}

	// Return valid environment
	return &Environment{DatabaseUrl, AuthHeader}, nil
}
