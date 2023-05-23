package logger

import (
	"errors"

	"github.com/joho/godotenv"
)

type Environment struct {
	DatabaseUrl string
}

func ValidateEnvironment() (*Environment, error) {
	// Parse .env
	var env map[string]string
	env, err := godotenv.Read()
	if err != nil {
		return &Environment{}, errors.New("failed to read env")
	}

	// Return valid environment
	return &Environment{DatabaseUrl: env["DATABASE_URL"]}, nil
}
