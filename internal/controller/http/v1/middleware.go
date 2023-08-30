package v1

import (
	"dynamic-user-segmentation/internal/service"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

const (
	prefix = "Bearer"
	uId    = "userId"
)

type AuthMiddleware struct {
	authService service.Auth
}

func (a *AuthMiddleware) CheckUser(handler echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		token, ok := checkToken(ctx.Request())
		if !ok {
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

//type TTLMiddleware struct {
//
//}
//
//func (ttl *TTLMiddleware) CheckTTL(handler echo.HandlerFunc) echo.HandlerFunc {
//	return func(ctx echo.Context) error {
//
//	}
//}

//type AutomaticAddMiddleware struct {
//}
//
//func (aa *AutomaticAddMiddleware) AutomaticAdd(handler echo.HandlerFunc) echo.HandlerFunc {
//	return func(ctx echo.Context) error {
//
//	}
//}
//
//func checkCount(count int, percent float64) bool {
//	number := int(math.Round(1 / percent))
//	if count%number == 0 {
//		return true
//	}
//	return false
//}
