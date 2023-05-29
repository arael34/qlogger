package logger

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type Environment struct {
	DatabaseUrl  string
	DatabaseName string
	AuthHeader   string
}

func ValidateEnvironment() (*Environment, error) {
	// In production, the environment variables
	// will be set in the OS env.
	if os.Getenv("ENV") == "PROD" {
		DatabaseUrl := os.Getenv("DATABASE_URL")
		DatabaseName := os.Getenv("DATABASE_NAME")
		AuthHeader := os.Getenv("AUTH_HEADER")

		if DatabaseUrl == "" || AuthHeader == "" {
			return nil, errors.New("failed to parse production env")
		}

		return &Environment{DatabaseUrl, DatabaseName, AuthHeader}, nil
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

	if DatabaseUrl == "" || DatabaseName == "" || AuthHeader == "" {
		return nil, errors.New("failed to parse local env")
	}

	return &Environment{DatabaseUrl, DatabaseName, AuthHeader}, nil
}
