package main

import (
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

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
	AuthHeader *string
	Database   *mongo.Collection
}

func NewQLogger(authHeader *string, database *mongo.Collection) *QLogger {
	return &QLogger{authHeader, database}
}

/*
 * Schema for a single log entry.
 */
type LogSchema struct {
	TimeWritten time.Time `json:"time" bson:"time"`
	Category    string    `json:"category" bson:"category"`
	Severity    Level     `json:"severity" bson:"severity"`
	Message     string    `json:"message" bson:"message"`
}
