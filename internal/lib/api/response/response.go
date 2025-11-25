package response

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

const (
	StatusOK    = "ok"
	StatusError = "error"
)

func OKResponse() Response {
	return Response{
		Status: StatusOK,
	}
}

func ErrorResponse(errMsg string) Response {
	return Response{
		Status: StatusError,
		Error:  errMsg,
	}
}

func ValidationError(errs validator.ValidationErrors) Response {
	var errMsg []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsg = append(errMsg, fmt.Sprintf("%s is a required field", err.Field()))
		case "url":
			errMsg = append(errMsg, fmt.Sprintf("%s is not a valid URL", err.Field()))
		default:
			errMsg = append(errMsg, fmt.Sprintf("%s is not valid", err.Field()))
		}
	}

	return Response{
		Status: StatusError,
		Error:  strings.Join(errMsg, ", "),
	}
}
