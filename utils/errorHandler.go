package utils

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
)

type ErrorResponseType struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// ErrorHandler struct for custom error handling

type ErrorHandler struct {
	StatusCode int
	Message    string
}

// NewErrorHandler function to create a new error instance

func NewErrorHandler(message string, statusCode int) *ErrorHandler {
	return &ErrorHandler{
		StatusCode: statusCode,
		Message:    message,
	}
}

// Implement the error interface for ErrorHandler

func (e *ErrorHandler) Error() string {
	return fmt.Sprintf("%s (Status Code: %d)", e.Message, e.StatusCode)
}

// To Show Log Error

func LogError(r *http.Request, err error) {

	logrus.Errorf("Error Received: %s %s %s", err, r.Method, r.URL.Path)
}

// ErrorResponse sends a JSON-formatted error response with the specified status code and message

func ErrorResponse(w http.ResponseWriter, r *http.Request, statusCode int, message string) {

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(statusCode)

	errResponseToSend := ErrorResponseType{
		Success: false,
		Message: message,
	}

	jsonData, err := json.Marshal(errResponseToSend)

	if err != nil {
		// If unable to marshal, log the error and send a generic error response
		LogError(r, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return

	}

	// Write JSON response to the ResponseWriter

	_, err = w.Write(jsonData)

	if err != nil {
		// If unable to write response, log the error
		LogError(r, err)
	}

}
