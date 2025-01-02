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
	stackTrace string
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
		datas = append(datas, string(ecRaw), d.stackTrace)
	}
	return fmt.Sprintf(format, datas...)
}

func (d *APIError) Inherit(err *APIError) *APIError {
	if err != nil {
		d.ErrorCause = err.ErrorCause
		d.stackTrace = err.stackTrace
	}
	return d
}

func New(status int, message string, args []string) *APIError {
	err := APIError{
		Status:  status,
		Message: message,
	}
	length := len(args)
	if length > 0 {
		err.ErrorCause = args[0]
		err.stackTrace = string(debug.Stack())
	}
	return &err
}

func Processing(message string, args ...string) *APIError {
	return New(http.StatusProcessing, message, args)
}

func Ok(message string, args ...string) *APIError {
	return New(http.StatusOK, message, args)
}

func Accepted(message string, args ...string) *APIError {
	return New(http.StatusAccepted, message, args)
}

func NoContent(message string, args ...string) *APIError {
	return New(http.StatusNoContent, message, args)
}

func TemporaryRedirect(message string, args ...string) *APIError {
	return New(http.StatusTemporaryRedirect, message, args)
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

func PreconditionFailed(message string, args ...string) *APIError {
	return New(http.StatusPreconditionFailed, message, args)
}

func RequestEntityTooLarge(message string, args ...string) *APIError {
	return New(http.StatusRequestEntityTooLarge, message, args)
}

func UnsupportedMediaType(message string, args ...string) *APIError {
	return New(http.StatusUnsupportedMediaType, message, args)
}

func UnprocessableEntity(message string, args ...string) *APIError {
	return New(http.StatusUnprocessableEntity, message, args)
}

func TooManyRequest(message string, args ...string) *APIError {
	return New(http.StatusTooManyRequests, message, args)
}

func InternalServerError(message string, args ...string) *APIError {
	return New(http.StatusInternalServerError, message, args)
}

func NotImplemented(message string, args ...string) *APIError {
	return New(http.StatusNotImplemented, message, args)
}

func ServiceUnavailable(message string, args ...string) *APIError {
	return New(http.StatusServiceUnavailable, message, args)
}
