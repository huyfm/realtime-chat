package rtc

import "errors"

const (
	ECONFLICT = iota + 1
	EINTERNAL
	EINVALID
	ENOTFOUND
	EUNAUTHORIZED
)

type Error struct {
	Code    int    // internal code
	Message string // client-oriented message
	Cause   error  // internal cause
}

func (e Error) Error() string {
	return e.Message
}

func ErrorCode(err error) int {
	if err == nil {
		return 0
	}
	var e *Error
	if errors.As(err, &e) {
		return e.Code
	}
	return EINTERNAL
}

func ErrorMsg(err error) string {
	if err == nil {
		return ""
	}
	var e *Error
	if errors.As(err, &e) {
		return e.Message
	}
	return "Internal error"
}

func Errorf(code int, msg string) *Error {
	return &Error{Code: code, Message: msg}
}
