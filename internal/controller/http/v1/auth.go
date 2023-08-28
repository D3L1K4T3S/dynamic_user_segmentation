package v1

import (
	"dynamic-user-segmentation/internal/service"
	"dynamic-user-segmentation/internal/service/dto"
	"errors"
	"github.com/labstack/echo"
	"net/http"
)

type signIn struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type signUp struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type authRoutes struct {
	authService service.Auth
}

func newAuthRoutes(group *echo.Group, authService service.Auth) {
	routes := &authRoutes{
		authService: authService,
	}

	group.POST("/sign-in", routes.signIn)
	group.POST("/sing-up", routes.signUp)
}

// @Description Sing In for users
// @Accept json
// @Produce json
// @Success 200
// @Failure 400
// @Failure 500
// @Router /auth/sign-in [post]
func (ar *authRoutes) signIn(ctx echo.Context) error {
	var input signIn

	if err := ctx.Bind(&input); err != nil {
		ErrResponse(ctx, http.StatusBadRequest, ErrInvalidRequest.Error())
		return err
	}

	token, err := ar.authService.GenerateToken(ctx.Request().Context(), dto.AuthUser{
		Username: input.Username,
		Password: input.Password,
	})

	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			ErrResponse(ctx, http.StatusBadRequest, ErrInvalidDataUser.Error())
			return err
		}
		ErrResponse(ctx, http.StatusInternalServerError, ErrInternalServer.Error())
		return err
	}

	return ctx.JSON(http.StatusOK, struct {
		Token string `json:"token"`
	}{token})
}

// @Description Sing Up for users
// @Accept json
// @Produce json
// @Success 201
// @Failure 400
// @Failure 500
// @Router /auth/sign-up [post]
func (ar *authRoutes) signUp(ctx echo.Context) error {
	var input signUp

	if err := ctx.Bind(&input); err != nil {
		ErrResponse(ctx, http.StatusBadRequest, ErrInvalidRequest.Error())
		return err
	}

	id, err := ar.authService.CreateUser(ctx.Request().Context(), dto.AuthUser{
		Username: input.Username,
		Password: input.Password,
	})

	if err != nil {
		if errors.Is(err, service.ErrUserAlreadyExists) {
			ErrResponse(ctx, http.StatusBadRequest, service.ErrUserAlreadyExists.Error())
			return err
		}
		ErrResponse(ctx, http.StatusInternalServerError, ErrInternalServer.Error())
		return err
	}

	return ctx.JSON(http.StatusCreated, struct {
		Id int `json:"id"`
	}{id})
}
