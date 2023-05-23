package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

/*
 * Simple int alias for readability.
 * 0 - DEBUG
 * 1 - WARN
 * 2 - ERROR
 */
type Level int

const (
	DEBUG Level = iota
	WARN
	ERROR
)

// This probably isn't needed, I'm keeping it just in case
type QLogger struct {
	DatabaseUrl *string
}

func NewQLogger(_DatabaseUrl *string) *QLogger {
	return &QLogger{DatabaseUrl: _DatabaseUrl}
}

/*
 * Schema for a single log entry.
 */
type LogSchema struct {
	Time    string // this needs to be changed to datetime
	Message string
	Level   Level
}

/*
 * Handler to write a log into the database.
 *
 * Expected body:
 *   Time: string (or datetime)
 *   Message: string
 *   Level: int
 */
func (s *QLogger) WriteLog(w http.ResponseWriter, r *http.Request) {
	var log LogSchema

	// Limit body size to 1MB and disallow unknown JSON fields
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1048576))
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&log)
	if err != nil {
		// TODO: better error handling
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Received message %v with level %v", log.Message, log.Level)
}

/*
 * Handler to read all logs. No body is necessary.
 *
 * Return body:
 *   Array<{
 *     Time: string (or datetime)
 *     Message: string
 *     Level: int
 *   }>
 */
func (s *QLogger) ReadLog(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("reading....")

	// Set response header to JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	// TODO this has to fetch from db
	var logs []LogSchema
	err := json.NewEncoder(w).Encode(logs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
