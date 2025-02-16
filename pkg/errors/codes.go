package errors

const (
	ENOTFOUND     = "ENOTFOUND"     // Resource not found
	EINVALID      = "EINVALID"      // Invalid input/validation errors
	ECONFLICT     = "ECONFLICT"     // Resource conflicts/already exists
	EINTERNAL     = "EINTERNAL"     // Internal server errors
	EUNAUTHORIZED = "EUNAUTHORIZED" // Authentication errors
	EFORBIDDEN    = "EFORBIDDEN"    // Authorization errors
	EBADREQUEST   = "EBADREQUEST"   // Malformed requests
	ETIMEOUT      = "ETIMEOUT"      // Operation timeouts
	ECACHE        = "ECACHE"        // Cache-related errors
	EVALIDATION   = "EVALIDATION"   // Domain validation errors
	EREPOSITORY   = "EREPOSITORY"   // Repository operation errors
)
