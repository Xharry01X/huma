package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"

	_ "github.com/danielgtaylor/huma/v2/formats/cbor" // Importing CBOR format support
)

// MyError is a custom error type that implements the huma.StatusError interface
type MyError struct {
	status  int      // HTTP status code
	Message string   `json:"message"`     // Error message
	Details []string `json:"details,omitempty"` // Optional details
}

// Error returns the error message
func (e *MyError) Error() string {
	return e.Message
}

// GetStatus returns the HTTP status code
func (e *MyError) GetStatus() int {
	return e.status
}

func main() {
	// Custom error handler function for huma
	huma.NewError = func(status int, message string, errs ...error) huma.StatusError {
		details := make([]string, len(errs))
		for i, err := range errs {
			details[i] = err.Error()
		}
		return &MyError{
			status:  status,
			Message: message,
			Details: details,
		}
	}

	// Initialize Chi router
	router := chi.NewMux()

	// Initialize humachi adapter with default configuration
	api := humachi.New(router, huma.DefaultConfig("My API", "1.0.0"))

	// Registering an operation with Huma
	huma.Register(api, huma.Operation{
		OperationID: "get-error",         // Unique operation ID
		Method:      http.MethodGet,      // HTTP method
		Path:        "/error",            // Endpoint path
	}, func(ctx context.Context, i *struct{}) (*struct{}, error) {
		// Example operation handler that returns a 404 Not Found error
		return nil, huma.Error404NotFound("not found", fmt.Errorf("some-other-error"))
	})

	// Start HTTP server on port 8888
	http.ListenAndServe(":8888", router)
}
