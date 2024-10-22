package apierrors

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/AAguilar0x0/bapp/pkg/assert"
	"runtime/debug"
)

type APIError struct {
	Status     int
	Message    string
	ErrorCause string
}

func (d *APIError) Error() string {
	msgRaw, err := json.Marshal(d.Message)
	assert.NoError(err, "marshalling Message", "fault", "Marshal")
	ecRaw, err := json.Marshal(d.ErrorCause)
	assert.NoError(err, "marshalling ErrorCause", "fault", "Marshal")
	return fmt.Sprintf("status=%d message=%s error=%s\nstack: %s", d.Status, string(msgRaw), string(ecRaw), string(debug.Stack()))
}

func New(status int, message string, args []string) *APIError {
	err := APIError{
		Status:  status,
		Message: message,
	}
	length := len(args)
	if length > 0 {
		err.ErrorCause = args[0]
	}
	return &err
}

func BadRequest(message string, args ...string) *APIError {
	return New(http.StatusBadRequest, message, args)
}

func Unauthorized(message string, args ...string) *APIError {
	return New(http.StatusUnauthorized, message, args)
}

func InternalServerError(message string, args ...string) *APIError {
	return New(http.StatusInternalServerError, message, args)
}
