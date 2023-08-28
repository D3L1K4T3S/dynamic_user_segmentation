package v1

import (
	"errors"
	"github.com/labstack/echo"
)

var (
	ErrInternalServer  = errors.New("internal error")
	ErrInvalidRequest  = errors.New("invalid request body")
	ErrInvalidDataUser = errors.New("invalid username or password")
)

func ErrResponse(ctx echo.Context, status int, message string) {
	err := errors.New(message)
	var HTTPError *echo.HTTPError
	ok := errors.As(err, &HTTPError)
	if !ok {
		report := echo.NewHTTPError(status, err.Error())
		_ = ctx.JSON(status, report)
	}
	ctx.Error(ErrInternalServer)
}
