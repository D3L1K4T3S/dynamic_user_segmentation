package respository_errors

import "errors"

var (
	ErrAlreadyExists = errors.New("record already exists")
	ErrNotFound      = errors.New("record not found")
	ErrCannotCreate  = errors.New("can't create record")
	ErrCannotAdd     = errors.New("can't add a record")
)
