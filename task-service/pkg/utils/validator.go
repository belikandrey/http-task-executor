package utils

import (
	"context"
	"github.com/go-playground/validator/v10"
	"http-task-executor/task-service/pkg/errors/general/validation"
	"net/http"
	"strings"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func ValidateStruct(ctx context.Context, data interface{}) error {
	return validate.StructCtx(ctx, data)
}

func ValidateHttpMethod(method string) validation.TaskValidationError {
	method = strings.ToUpper(method)

	if method != http.MethodGet && method != http.MethodHead && method != http.MethodPost &&
		method != http.MethodPut && method != http.MethodPatch && method != http.MethodDelete &&
		method != http.MethodConnect && method != http.MethodOptions && method != http.MethodTrace {
		return validation.CustomFiledError{Fld: "Method", Msg: "invalid http method", Tag: "http-method"}
	}
	return nil
}
