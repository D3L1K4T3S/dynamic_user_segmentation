package service

import "errors"

var (
	ErrCannotCreateToken       = errors.New("can't create a token")
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")
	ErrCannotParseToken        = errors.New("can't parse a token")
)
