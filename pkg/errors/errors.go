package errors

type AppError struct {
	Err     error
	Message string
	Code    string
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
