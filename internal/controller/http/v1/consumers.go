package v1

import (
	"dynamic-user-segmentation/internal/service"
	"dynamic-user-segmentation/internal/service/dto"
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
)

type consumersRoutes struct {
	consumersService service.Consumers
}

func newConsumersRoutes(group *echo.Group, consumers service.Consumers) {
	routes := &consumersRoutes{
		consumers,
	}

	group.POST("/create", routes.create)
	group.PUT("/add", routes.add)
	group.GET("/get", routes.get)
	group.PATCH("/update", routes.update)
	group.DELETE("/delete", routes.delete)
}

// @Description Create user
// @Tags consumers
// @Accept json
// @Produce json
// @Param input body dto.ConsumerRequest true "input"
// @Success 201
// @Failure 400
// @Failure 500
// @Security JWT
// @Router /api/v1/consumers/create [post]
func (cr *consumersRoutes) create(ctx echo.Context) error {
	var input dto.ConsumerRequest

	if err := ctx.Bind(&input); err != nil {
		ErrResponse(ctx, http.StatusBadRequest, ErrInvalidRequest.Error())
		return err
	}

	id, err := cr.consumersService.CreateConsumer(ctx.Request().Context(), input)
	if err != nil {
		if errors.Is(err, service.ErrUserAlreadyExists) {
			ErrResponse(ctx, http.StatusBadRequest, service.ErrUserAlreadyExists.Error())
			return err
		}
		ErrResponse(ctx, http.StatusInternalServerError, ErrInternalServer.Error())
		return err
	}
	return ctx.JSON(http.StatusCreated, struct {
		Id []int `json:"id"`
	}{id})
}

// @Description Add segments to consumer
// @Tags consumers
// @Accept json
// @Produce json
// @Param input body dto.ConsumerRequest true "input"
// @Success 201
// @Failure 400
// @Failure 500
// @Security JWT
// @Router /api/v1/consumers/add [put]
func (cr *consumersRoutes) add(ctx echo.Context) error {
	var input dto.ConsumerRequest

	if err := ctx.Bind(&input); err != nil {
		ErrResponse(ctx, http.StatusBadRequest, ErrInvalidRequest.Error())
		return err
	}

	id, err := cr.consumersService.AddSegmentsToConsumer(ctx.Request().Context(), input)

	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			ErrResponse(ctx, http.StatusBadRequest, service.ErrUserNotFound.Error())
			return err
		}
		ErrResponse(ctx, http.StatusInternalServerError, ErrInternalServer.Error())
		return err
	}

	return ctx.JSON(http.StatusOK, struct {
		Id []int `json:"id"`
	}{id})
}

// @Description Get all segments consumer
// @Tags consumers
// @Accept json
// @Produce json
// @Param input body dto.ConsumerId true "input"
// @Success 201
// @Failure 400
// @Failure 500
// @Security JWT
// @Router /api/v1/consumers/get [get]
func (cr *consumersRoutes) get(ctx echo.Context) error {
	var input dto.ConsumerId

	if err := ctx.Bind(&input); err != nil {
		ErrResponse(ctx, http.StatusBadRequest, ErrInvalidRequest.Error())
		return err
	}

	res, err := cr.consumersService.GetConsumerSegments(ctx.Request().Context(), input)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			ErrResponse(ctx, http.StatusBadRequest, ErrInvalidRequest.Error())
			return err
		}
		ErrResponse(ctx, http.StatusInternalServerError, ErrInternalServer.Error())
		return err
	}

	return ctx.JSON(http.StatusOK, res)
}

// @Description Update ttl for segment consumer
// @Tags consumers
// @Accept json
// @Produce json
// @Param input body dto.ConsumerRequest true "input"
// @Success 200
// @Failure 400
// @Failure 500
// @Security JWT
// @Router /api/v1/consumers/update [patch]
func (cr *consumersRoutes) update(ctx echo.Context) error {
	var input dto.ConsumerRequest

	if err := ctx.Bind(&input); err != nil {
		ErrResponse(ctx, http.StatusBadRequest, ErrInvalidRequest.Error())
		return err
	}

	err := cr.consumersService.UpdateSegmentsTTL(ctx.Request().Context(), input)

	if err != nil {
		ErrResponse(ctx, http.StatusInternalServerError, ErrInternalServer.Error())
		return err
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
	})
}

// @Description Delete user segment
// @Tags consumers
// @Accept json
// @Produce json
// @Param input body dto.ConsumerRequest true "input"
// @Success 200
// @Failure 400
// @Failure 500
// @Security JWT
// @Router /api/v1/consumers/delete [delete]
func (cr *consumersRoutes) delete(ctx echo.Context) error {
	var input dto.ConsumerRequestDelete

	if err := ctx.Bind(&input); err != nil {
		ErrResponse(ctx, http.StatusBadRequest, ErrInvalidRequest.Error())
		return err
	}

	err := cr.consumersService.DeleteSegmentsFromConsumer(ctx.Request().Context(), input)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			ErrResponse(ctx, http.StatusBadRequest, service.ErrUserNotFound.Error())
			return err
		}
		if errors.Is(err, service.ErrSegmentNotFound) {
			ErrResponse(ctx, http.StatusBadRequest, service.ErrSegmentNotFound.Error())
			return err
		}

		ErrResponse(ctx, http.StatusInternalServerError, ErrInternalServer.Error())
		return err
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
	})
}
