package apierrors

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/AAguilar0x0/bapp/pkg/assert"
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
	return fmt.Sprintf("cause=%s status=%d message=%s", string(ecRaw), d.Status, string(msgRaw))
}

func new(status int, message string, args []string) *APIError {
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

func InternalServerError(message string, args ...string) *APIError {
	return new(http.StatusInternalServerError, message, args)
}

func Unauthorized(message string, args ...string) *APIError {
	return new(http.StatusUnauthorized, message, args)
}
