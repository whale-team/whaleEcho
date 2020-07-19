package wserrors

import (
	"strings"
)

var (
	// ErrNotFound represent resource not found error
	ErrNotFound     = New(NotFound, "resource not found")
	ErrNotAuth      = New(NotAuth, "websocket message is unauthorized")
	ErrSysBusy      = New(SysBusy, "webscoket message timeout")
	ErrInternal     = New(Internal, "The server encountered an internal error. Please notify admin")
	ErrInputInvalid = New(InputInvalid, "websocket message format is invalid")
)

type ErrStatus int64

const (
	NotAuth      ErrStatus = 3
	NotFound     ErrStatus = 4
	SysBusy      ErrStatus = 5
	Internal     ErrStatus = 6
	InputInvalid ErrStatus = 7
)

var StatusMap = map[ErrStatus]string{
	NotAuth:      "not auth",
	NotFound:     "not found",
	SysBusy:      "system busy",
	Internal:     "internal error",
	InputInvalid: "input invalid",
}

func New(status ErrStatus, message string) error {
	return WithStack(&WsError{
		Status:  status,
		Message: message,
	})
}

type WsError struct {
	Status  ErrStatus
	Message string
}

func (err *WsError) Error() string {
	var b strings.Builder
	_, _ = b.WriteRune('[')
	_, _ = b.WriteString(StatusMap[err.Status])
	_, _ = b.WriteRune(']')
	_, _ = b.WriteRune(' ')
	_, _ = b.WriteString(err.Message)
	return b.String()
}
