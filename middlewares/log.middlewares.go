package middlewares

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

// Custom ResponseWriter to capture the status code
type statusRecorder struct {
	http.ResponseWriter
	StatusCode int
}

// Overriding WriteHeader to capture the status code
func (rec *statusRecorder) WriteHeader(code int) {
	rec.StatusCode = code
	rec.ResponseWriter.WriteHeader(code)
}

// LogMiddleware logs request details and response status
func LogMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Create a statusRecorder to capture the status code
		rec := &statusRecorder{
			ResponseWriter: w,
			StatusCode:     http.StatusOK, // Default status code to 200
		}

		// Log the incoming request
		logrus.Infof("Request Received: %s %s", r.Method, r.URL.Path)

		// Call the next handler in the chain
		next(rec, r)

		// Log the response status code after the handler has been executed
		logrus.Infof("Response Status: %d %s %s", rec.StatusCode, r.Method, r.URL.Path)
	}
}
