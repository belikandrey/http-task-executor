package http

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"http-task-executor/task-service/pkg/errors/general/validation"
	"net/http"
	"strconv"
	"strings"
)

var (
	// ErrBadRequest - bad request error
	ErrBadRequest = errors.New("bad request")
	// ErrNotFound - not found error
	ErrNotFound = errors.New("not Found")
	// ErrRequestTimeoutError - request timeout error
	ErrRequestTimeoutError = errors.New("request Timeout")
	// ErrInternalServerError - internal service error
	ErrInternalServerError = errors.New("internal Server Error")
)

// RestErr represents REST error.
type RestErr interface {
	Status() int
	Error() string
	Causes() interface{}
}

// RestError represents REST error and implements RestErr.
type RestError struct {
	ErrStatus int         `json:"status,omitempty"`
	ErrError  string      `json:"error,omitempty"`
	ErrCauses interface{} `json:"-"`
}

// Error returns formatted error message
func (e RestError) Error() string {
	return fmt.Sprintf("status: %d - errors: %s - causes: %v", e.ErrStatus, e.ErrError, e.ErrCauses)
}

// Status returns http status code
func (e RestError) Status() int {
	return e.ErrStatus
}

// Causes returns causes
func (e RestError) Causes() interface{} {
	return e.ErrCauses
}

// NewBadRequestError creates new bad request error
func NewBadRequestError(causes interface{}) RestErr {
	return RestError{
		ErrStatus: http.StatusBadRequest,
		ErrError:  ErrBadRequest.Error(),
		ErrCauses: causes,
	}
}

// NewValidationError creates new validation error
func NewValidationError(errs []validation.TaskValidationError) RestErr {

	var errMsgs []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s: is a required field", err.Field()))
		case "url":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s: is not a valid URL", err.Field()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field %s: %s", err.Field(), err.Error()))
		}
	}

	return RestError{
		ErrStatus: http.StatusBadRequest,
		ErrError:  strings.Join(errMsgs, ", "),
	}
}

// ErrorResponse returns status code and response object
func ErrorResponse(err error) (int, interface{}) {
	return ParseErrors(err).Status(), ParseErrors(err)
}

// NewRestError creates error
func NewRestError(status int, err string, causes interface{}) RestErr {
	return RestError{
		ErrStatus: status,
		ErrError:  err,
		ErrCauses: causes,
	}
}

// NewInternalServerError creates internal server error
func NewInternalServerError(causes interface{}) RestErr {
	result := RestError{
		ErrStatus: http.StatusInternalServerError,
		ErrError:  ErrInternalServerError.Error(),
		ErrCauses: causes,
	}
	return result
}

// ParseErrors parses error and returns RestErr based on error type.
func ParseErrors(err error) RestErr {
	var unmarshalTypeError *json.UnmarshalTypeError
	var jsonSyntaxType *json.SyntaxError
	var strconvNumError *strconv.NumError
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return NewRestError(http.StatusNotFound, ErrNotFound.Error(), err)
	case errors.Is(err, context.DeadlineExceeded):
		return NewRestError(http.StatusRequestTimeout, ErrRequestTimeoutError.Error(), err)
	case errors.As(err, &unmarshalTypeError):
		return NewRestError(http.StatusBadRequest, ErrBadRequest.Error(), err)
	case errors.As(err, &jsonSyntaxType):
		return NewRestError(http.StatusBadRequest, ErrBadRequest.Error(), err)
	case errors.As(err, &strconvNumError):
		return NewRestError(http.StatusBadRequest, ErrBadRequest.Error(), err)
	default:
		if restErr, ok := err.(RestErr); ok {
			return restErr
		}
		return NewInternalServerError(err)
	}
}
