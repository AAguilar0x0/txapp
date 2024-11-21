package apierrors

import (
	"encoding/json"
	"fmt"
	"net/http"

	"runtime/debug"

	"github.com/AAguilar0x0/txapp/core/pkg/assert"
)

type APIError struct {
	Status     int
	Message    string
	ErrorCause string
}

func (d *APIError) Error() string {
	format := "status=%d message=%s"
	msgRaw, err := json.Marshal(d.Message)
	assert.NoError(err, "marshalling Message", "fault", "Marshal")
	datas := []interface{}{
		d.Status,
		string(msgRaw),
	}
	if d.ErrorCause != "" {
		ecRaw, err := json.Marshal(d.ErrorCause)
		assert.NoError(err, "marshalling ErrorCause", "fault", "Marshal")
		format += " error=%s\nstack: %s"
		datas = append(datas, string(ecRaw), string(debug.Stack()))
	}
	return fmt.Sprintf(format, datas...)
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

func NotFound(message string, args ...string) *APIError {
	return New(http.StatusNotFound, message, args)
}

func Conflict(message string, args ...string) *APIError {
	return New(http.StatusConflict, message, args)
}

func Unauthorized(message string, args ...string) *APIError {
	return New(http.StatusUnauthorized, message, args)
}

func Forbidden(message string, args ...string) *APIError {
	return New(http.StatusForbidden, message, args)
}

func RequestTimeout(message string, args ...string) *APIError {
	return New(http.StatusRequestTimeout, message, args)
}

func InternalServerError(message string, args ...string) *APIError {
	return New(http.StatusInternalServerError, message, args)
}

func ServiceUnavailable(message string, args ...string) *APIError {
	return New(http.StatusServiceUnavailable, message, args)
}
