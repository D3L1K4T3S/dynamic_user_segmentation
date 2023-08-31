package v1

import (
	"dynamic-user-segmentation/internal/service"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"strings"
)

const (
	prefix = "Bearer "
	uId    = "userId"
)

type AuthMiddleware struct {
	authService service.Auth
}

func (a *AuthMiddleware) CheckUser(handler echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		token, ok := checkToken(ctx.Request())
		if !ok {
			log.Println("Use bearerToken: ", ErrInvalidAuth)
			ErrResponse(ctx, http.StatusUnauthorized, ErrInvalidAuth.Error())
			return nil
		}

		userId, err := a.authService.ParseToken(token)
		if err != nil {
			ErrResponse(ctx, http.StatusUnauthorized, service.ErrCannotParseToken.Error())
			return err
		}

		ctx.Set(uId, userId)
		return handler(ctx)
	}
}

func checkToken(r *http.Request) (string, bool) {
	header := r.Header.Get(echo.HeaderAuthorization)
	if header == "" {
		return "", false
	}

	if strings.EqualFold(header[:len(prefix)], prefix) && len(header) > len(prefix) {
		return header[len(prefix):], true
	}
	return "", false
}
