package v1

import (
	"dynamic-user-segmentation/internal/service"
	"dynamic-user-segmentation/internal/service/dto"
	"errors"
	"github.com/labstack/echo"
	"net/http"
)

type segmentsRoutes struct {
	segmentService service.Segments
}

func newSegmentRoutes(group *echo.Group, segmentsService service.Segments) {
	routes := &segmentsRoutes{
		segmentService: segmentsService,
	}

	group.POST("/create", routes.create)
	group.PATCH("/update", routes.update)
	group.DELETE("/delete", routes.delete)
}

type segmentInput struct {
	Name    string
	Percent float64
}

type segmentDelete struct {
	Name string
}

// @Description Create segment
// @Accept json
// @Produce json
// @Success 201
// @Failure 400
// @Failure 500
// @Router /api/v1/segments/create [post]
func (sr *segmentsRoutes) create(ctx echo.Context) error {
	var input segmentInput

	if err := ctx.Bind(&input); err != nil {
		ErrResponse(ctx, http.StatusBadRequest, ErrInvalidRequest.Error())
		return err
	}

	id, err := sr.segmentService.CreateSegment(ctx.Request().Context(), dto.SegmentsRequest{
		Name:    input.Name,
		Percent: input.Percent,
	})

	if err != nil {
		if errors.Is(err, service.ErrSegmentAlreadyExists) {
			ErrResponse(ctx, http.StatusBadRequest, service.ErrSegmentAlreadyExists.Error())
			return err
		}
		ErrResponse(ctx, http.StatusInternalServerError, ErrInternalServer.Error())
		return err
	}

	return ctx.JSON(http.StatusCreated, struct {
		Id int `json:"id"`
	}{id})
}

// @Description Update percent in segment
// @Accept json
// @Produce json
// @Success 200
// @Failure 400
// @Failure 500
// @Router /api/v1/segments/update [patch]
func (sr *segmentsRoutes) update(ctx echo.Context) error {
	var input segmentInput

	if err := ctx.Bind(&input); err != nil {
		ErrResponse(ctx, http.StatusBadRequest, ErrInvalidRequest.Error())
		return err
	}

	err := sr.segmentService.UpdateSegment(ctx.Request().Context(), dto.SegmentsRequest{
		Name:    input.Name,
		Percent: input.Percent,
	})

	if err != nil {
		if errors.Is(err, service.ErrSegmentAlreadyExists) {
			ErrResponse(ctx, http.StatusBadRequest, service.ErrSegmentAlreadyExists.Error())
			return err
		}
		ErrResponse(ctx, http.StatusInternalServerError, ErrInternalServer.Error())
		return err
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
	})
}

// @Description Delete segment
// @Accept json
// @Produce json
// @Success 200
// @Failure 400
// @Failure 500
// @Router /api/v1/segments/delete [delete]
func (sr *segmentsRoutes) delete(ctx echo.Context) error {
	var input segmentDelete

	if err := ctx.Bind(&input); err != nil {
		ErrResponse(ctx, http.StatusBadRequest, ErrInvalidRequest.Error())
		return err
	}

	err := sr.segmentService.DeleteSegment(ctx.Request().Context(), dto.SegmentsRequest{
		Name: input.Name,
	})

	if err != nil {
		if errors.Is(err, service.ErrSegmentAlreadyExists) {
			ErrResponse(ctx, http.StatusBadRequest, service.ErrSegmentAlreadyExists.Error())
			return err
		}
		ErrResponse(ctx, http.StatusInternalServerError, ErrInternalServer.Error())
		return err
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
	})
}
