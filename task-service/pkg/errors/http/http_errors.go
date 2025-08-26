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
	ErrBadRequest          = errors.New("bad request")
	ErrNotFound            = errors.New("not Found")
	ErrRequestTimeoutError = errors.New("request Timeout")
	ErrInternalServerError = errors.New("internal Server Error")
)

type RestErr interface {
	Status() int
	Error() string
	Causes() interface{}
}

type RestError struct {
	ErrStatus int         `json:"status,omitempty"`
	ErrError  string      `json:"error,omitempty"`
	ErrCauses interface{} `json:"-"`
}

func (e RestError) Error() string {
	return fmt.Sprintf("status: %d - errors: %s - causes: %v", e.ErrStatus, e.ErrError, e.ErrCauses)
}

func (e RestError) Status() int {
	return e.ErrStatus
}

func (e RestError) Causes() interface{} {
	return e.ErrCauses
}

func NewBadRequestError(causes interface{}) RestErr {
	return RestError{
		ErrStatus: http.StatusBadRequest,
		ErrError:  ErrBadRequest.Error(),
		ErrCauses: causes,
	}
}

func NewValidationError(errs []validation.ValidationError) RestErr {

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

func ErrorResponse(err error) (int, interface{}) {
	return ParseErrors(err).Status(), ParseErrors(err)
}

func NewRestError(status int, err string, causes interface{}) RestErr {
	return RestError{
		ErrStatus: status,
		ErrError:  err,
		ErrCauses: causes,
	}
}

func NewInternalServerError(causes interface{}) RestErr {
	result := RestError{
		ErrStatus: http.StatusInternalServerError,
		ErrError:  ErrInternalServerError.Error(),
		ErrCauses: causes,
	}
	return result
}

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
