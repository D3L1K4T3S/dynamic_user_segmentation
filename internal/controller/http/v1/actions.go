package v1

import (
	"dynamic-user-segmentation/internal/service"
	"dynamic-user-segmentation/internal/service/dto"
	"errors"
	"github.com/labstack/echo"
	"net/http"
)

type actionsRoutes struct {
	actionsService service.Actions
}

func newActionsRoutes(group *echo.Group, actionsService service.Actions) {
	routes := &actionsRoutes{
		actionsService: actionsService,
	}

	group.POST("/create", routes.create)
	group.DELETE("/delete", routes.delete)
}

type actionInput struct {
	Name string `json:"name"`
}

// @Description Create action
// @Accept json
// @Produce json
// @Success 201
// @Failure 400
// @Failure 500
// @Router /api/v1/actions/create [post]
func (ar *actionsRoutes) create(ctx echo.Context) error {
	var input actionInput

	if err := ctx.Bind(&input); err != nil {
		ErrResponse(ctx, http.StatusBadRequest, ErrInvalidRequest.Error())
		return err
	}

	id, err := ar.actionsService.CreateAction(ctx.Request().Context(), dto.ActionsRequest{
		Name: input.Name,
	})

	if err != nil {
		if errors.Is(err, service.ErrActionAlreadyExists) {
			ErrResponse(ctx, http.StatusBadRequest, service.ErrUserAlreadyExists.Error())
			return err
		}
		ErrResponse(ctx, http.StatusInternalServerError, ErrInternalServer.Error())
		return err
	}
	return ctx.JSON(http.StatusOK, struct {
		Id int `json:"id"`
	}{id})
}

// @Description Delete action
// @Accept json
// @Produce json
// @Success 200
// @Failure 400
// @Failure 500
// @Router /api/v1/actions/delete [delete]
func (ar *actionsRoutes) delete(ctx echo.Context) error {
	var input actionInput

	if err := ctx.Bind(&input); err != nil {
		ErrResponse(ctx, http.StatusBadRequest, ErrInvalidRequest.Error())
		return err
	}

	err := ar.actionsService.DeleteAction(ctx.Request().Context(), dto.ActionsRequest{
		Name: input.Name,
	})
	if err != nil {
		if errors.Is(err, service.ErrActionNotFound) {
			ErrResponse(ctx, http.StatusBadRequest, service.ErrUserNotFound.Error())
			return err
		}
		ErrResponse(ctx, http.StatusInternalServerError, ErrInternalServer.Error())
		return err
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
	})
}
