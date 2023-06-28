package types

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
	Severity    int       `json:"severity" bson:"severity"`
	Category    string    `json:"category" bson:"category"`
}
