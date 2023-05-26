package logger

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
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

	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(15*time.Second))
	defer cancel()

	// Insert schema into database.
	_, err = logger.database.InsertOne(ctx, log)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

/*
 * Handler to read all logs. No body is necessary.
 * Expects Authorization header.
 *
 * Filter through an optional query parameter.
 *   ?severity={int}
 *   ?origin={string} <- case sensitive
 *
 * Return body:
 *   data: Array<{
 *     Time: datetime
 *     Message: string
 *     Severity: int
 *   }>
 */
func (logger *QLogger) ReadLogs(w http.ResponseWriter, r *http.Request) {
	// Authorize user.
	if r.Header.Get("Authorization") != *logger.authHeader {
		http.Error(w, "not authorized.", http.StatusUnauthorized)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "method not allowed.", http.StatusMethodNotAllowed)
	}

	query := r.URL.Query()
	filter := bson.M{} // bson.M{} applies no filter.

	// Apply filters
	if severity := query.Get("severity"); severity != "" {
		filter["severity"] = severity

		// Try to convert to an int
		convertedSeverity, err := strconv.Atoi(severity)
		if err == nil {
			filter["severity"] = convertedSeverity
		}
	}
	if origin := query.Get("origin"); origin != "" {
		filter["origin"] = origin
	}

	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(15*time.Second))
	defer cancel()

	cursor, err := logger.database.Find(ctx, filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Parse documents as schema.
	var logs []LogSchema
	err = cursor.All(ctx, &logs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(logs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
