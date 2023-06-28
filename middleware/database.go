package middleware

import "context"

type LogSchema struct {
	TimeWritten string `json:"time_written" bson:"time_written"`
	Message     string `json:"message" bson:"message"`
	Severity    int    `json:"severity" bson:"severity"`
	Category    string `json:"category" bson:"category"`
}

type Database interface {
	Write(context.Context, LogSchema) error
}
