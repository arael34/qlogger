package logger

import (
	"fmt"
	"net/http"
	"time"

	"github.com/jonasiwnl/qlogger/types"
)

// LogRoute middleware logs the request to the database.
func (q QLogger) LogRoute(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// DEBUG
		fmt.Println("reached log route.")

		q.database.Write(r.Context(), types.LogSchema{
			TimeWritten: time.Now().UTC(),
			Message:     "Request to " + r.URL.Path + " from " + r.RemoteAddr,
			Severity:    0,
			Category:    "api",
		})

		// DEBUG
		fmt.Println("wrote log.")

		next.ServeHTTP(w, r)
	})
}
