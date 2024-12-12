package http

import (
	"net/http"

	"github.com/huyfm/rtc"
)

type Session struct {
	UserID int `json:"userID,omitempty"`
	// OAuth2 random state to prevent csrf.
	State string `json:"state,omitempty"`
}

const sessionName = "session"

var httpCodes = map[int]int{
	rtc.ECONFLICT:     http.StatusConflict,
	rtc.EINTERNAL:     http.StatusInternalServerError,
	rtc.EINVALID:      http.StatusBadRequest,
	rtc.ENOTFOUND:     http.StatusNotFound,
	rtc.EUNAUTHORIZED: http.StatusUnauthorized,
}

// Error respones with error message and optionally log error.
func Error(w http.ResponseWriter, r *http.Request, err error) {
	code, msg := rtc.ErrorCode(err), rtc.ErrorMsg(err)
	if code == 0 {
		return
	}
	// Log internal error.
	if code == rtc.EINTERNAL {
		LogError(r, err)
	}
	w.WriteHeader(httpCodes[code])
	w.Write([]byte(msg))
}

func LogError(r *http.Request, err error) {
	rtc.Logger.Warn().
		Str("method", r.Method).
		Str("path", r.URL.Path).
		Err(err).
		Msg("")
}

func LogInfo(r *http.Request, msg string) {
	rtc.Logger.Info().
		Str("method", r.Method).
		Str("path", r.URL.Path).
		Msg(msg)
}

func LogOK(r *http.Request) {
	LogInfo(r, "OK")
}
