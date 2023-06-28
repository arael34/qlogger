package lib

import (
	"context"
	"time"
)

/*
 * Database interface for writing logs.
 */
type Database interface {
	Write(context.Context, LogSchema) error
}

/*
 * Schema for a single log entry.
 */
type LogSchema struct {
	TimeWritten time.Time `json:"time_written" bson:"time_written"`
	Message     string    `json:"message" bson:"message"`
	Severity    Level     `json:"severity" bson:"severity"`
	Category    string    `json:"category" bson:"category"`
}

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
