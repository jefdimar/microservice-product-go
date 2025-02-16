package errors

import (
	"fmt"
	"strings"
)

type AppError struct {
	Err     error
	Message string
	Code    string
}

// Implement the error interface
func (e *AppError) Error() string {
	var b strings.Builder

	if e.Code != "" {
		fmt.Fprintf(&b, "[%s] ", e.Code)
	}

	if e.Message != "" {
		fmt.Fprintf(&b, "%s", e.Message)
	}

	if e.Err != nil {
		if e.Message != "" {
			b.WriteString(": ")
		}
		fmt.Fprintf(&b, "%v", e.Err)
	}

	return b.String()
}

func StandardError(code string, err error) *AppError {
	var msg string
	switch code {
	case ENOTFOUND:
		msg = MsgNotFound
	case EINVALID:
		msg = MsgInvalidInput
	case ECONFLICT:
		msg = MsgAlreadyExists
	case EINTERNAL:
		msg = MsgInternalError
	case EUNAUTHORIZED:
		msg = MsgUnauthorized
	case EFORBIDDEN:
		msg = MsgForbidden
	case EBADREQUEST:
		msg = MsgBadRequest
	case ETIMEOUT:
		msg = MsgTimeout
	case ECACHE:
		msg = MsgCacheError
	case EVALIDATION:
		msg = MsgValidationFailed
	case EREPOSITORY:
		msg = MsgRepositoryError
	default:
		msg = MsgInternalError
	}

	return &AppError{
		Code:    code,
		Message: msg,
		Err:     err,
	}
}
