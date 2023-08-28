package service

import "errors"

var (
	ErrCannotCreateToken       = errors.New("can't create a token")
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")
	ErrCannotParseToken        = errors.New("can't parse a token")

	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")

	ErrActionAlreadyExists = errors.New("action already exists")
	ErrActionNotFound      = errors.New("action not found")

	ErrSegmentAlreadyExists = errors.New("segment already exists")
	ErrSegmentNotFound      = errors.New("segment not found")
)
