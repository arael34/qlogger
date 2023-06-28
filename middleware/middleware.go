package middleware

import (
	"fmt"
	"net/http"
	"time"
)

type QMiddleware struct {
	database Database
}

// LogRoute middleware logs the request to the database.
func (q QMiddleware) LogRoute(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// DEBUG
		fmt.Println("reached log route.")

		q.database.Write(r.Context(), LogSchema{
			TimeWritten: time.Now().UTC().Format(time.RFC3339),
			Message:     "Request to " + r.URL.Path + " from " + r.RemoteAddr,
			Severity:    0,
			Category:    "api",
		})

		// DEBUG
		fmt.Println("wrote log.")

		next.ServeHTTP(w, r)
	})
}
