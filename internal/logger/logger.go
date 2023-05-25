package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

/*
 * Simple alias for readability.
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

type QLogger struct {
	authHeader *string
	database   *mongo.Collection
}

func NewQLogger(authHeader *string, database *mongo.Collection) *QLogger {
	return &QLogger{authHeader, database}
}

/*
 * Schema for a single log entry.
 */
type LogSchema struct {
	TimeWritten time.Time `bson:"time"`
	Message     string    `bson:"message"`
	Origin      string    `bson:"origin"`
	Severity    Level     `bson:"severity"`
}

/*
 * Handler to write a log into the database.
 * Expects Authorization header.
 *
 * Expected body:
 *   Message: string
 *   Severity: int
 *   Origin: string
 */
func (logger *QLogger) WriteLog(w http.ResponseWriter, r *http.Request) {
	// Authorize user.
	if r.Header.Get("Authorization") != *logger.authHeader {
		http.Error(w, "not authorized.", http.StatusUnauthorized)
		return
	}
	if r.Method != "POST" {
		http.Error(w, "method not allowed.", http.StatusMethodNotAllowed)
	}

	fmt.Println("writing...")

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

	log.TimeWritten = time.Now().UTC()

	// Insert schema into database.
	_, err = logger.database.InsertOne(context.TODO(), log)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("wrote new log.")
}

/*
 * Handler to read all logs. No body is necessary.
 * Expects Authorization header.
 *
 * Return body:
 *   data: Array<{
 *     Time: datetime
 *     Message: string
 *     Level: int
 *   }>
 */
func (logger *QLogger) ReadLog(w http.ResponseWriter, r *http.Request) {
	// Authorize user.
	if r.Header.Get("Authorization") != *logger.authHeader {
		http.Error(w, "not authorized.", http.StatusUnauthorized)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "method not allowed.", http.StatusMethodNotAllowed)
	}

	fmt.Println("reading....")

	// To fetch every entry, set filter to bson.D{}.
	cursor, err := logger.database.Find(context.TODO(), bson.D{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Parse documents as schema.
	var logs []LogSchema
	err = cursor.All(context.TODO(), &logs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return parsed schema as JSON.
	err = json.NewEncoder(w).Encode(logs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println("read all logs.")
}
