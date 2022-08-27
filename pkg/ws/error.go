package ws

import "errors"

var (
	ErrServer               = errors.New("server error")
	ErrRemoteServerErr      = errors.New("remote server error")
	ErrInvalidCredential    = errors.New("invalid credential")
	ErrFailedToReadResponse = errors.New("failed to read response")
)
