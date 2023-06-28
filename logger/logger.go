package logger

import "github.com/jonasiwnl/qlogger/types"

/*
 * Simple alias for readability.
 * 0 - INFO
 * 1 - DEBUG
 * 2 - WARN
 * 3 - ERROR
 */
type Level int

const (
	INFO Level = iota
	DEBUG
	WARN
	ERROR
)

type QLogger struct {
	database       types.Database
	allowedOrigins *[]string
}

func NewQLogger(database types.Database, allowedOrigins *[]string) *QLogger {
	return &QLogger{database, allowedOrigins}
}
