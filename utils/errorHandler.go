package utils

import "fmt"

// ErrorHandler struct for custom error handling

type ErrorHandler struct {
	StatusCode int 
	Message string
}

// NewErrorHandler function to create a new error instance

func NewErrorHandler(message string, statusCode int) *ErrorHandler {
	return &ErrorHandler{
		StatusCode: statusCode ,
		Message: message,
		
	}
}

// Implement the error interface for ErrorHandler

func (e *ErrorHandler) Error() string {
	return fmt.Sprintf("%s (Status Code: %d)", e.Message, e.StatusCode)
}




