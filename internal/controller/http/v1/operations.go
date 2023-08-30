package v1

import (
	"dynamic-user-segmentation/internal/service"
	"dynamic-user-segmentation/internal/service/dto"
	"github.com/labstack/echo/v4"
	"net/http"
)

type operationRoutes struct {
	service.Operations
}

func newOperationsRoutes(group *echo.Group, operations service.Operations) {
	routes := &operationRoutes{
		operations,
	}

	group.GET("/", routes.getOperations)
	group.GET("/file", routes.getFileLink)

}

// @Description Get history operations json
// @Tags operations
// @Accept json
// @Produce json
// @Param input body dto.OperationsRequest true "input"
// @Success 201
// @Failure 400 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Security JWT
// @Router /api/v1/operations/ [get]
func (or *operationRoutes) getOperations(ctx echo.Context) error {
	var input dto.OperationsRequest

	if err := ctx.Bind(&input); err != nil {
		ErrResponse(ctx, http.StatusBadRequest, ErrInvalidRequest.Error())
		return err
	}

	operations, err := or.Operations.GetHistoryForPeriod(ctx.Request().Context(), input)

	if err != nil {
		ErrResponse(ctx, http.StatusInternalServerError, ErrInternalServer.Error())
		return err
	}

	return ctx.JSON(http.StatusOK, operations)

}

// @Description Get file
// @Tags operations
// @Accept json
// @Produce text/csv
// @Param input body dto.OperationsRequest true "input"
// @Success 201
// @Failure 400 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Security JWT
// @Router /api/v1/operations/file [get]
func (or *operationRoutes) getFileLink(ctx echo.Context) error {
	var input dto.OperationsRequest

	if err := ctx.Bind(&input); err != nil {
		ErrResponse(ctx, http.StatusBadRequest, ErrInvalidRequest.Error())
		return err
	}

	file, err := or.Operations.GetReportFile(ctx.Request().Context(), input)
	if err != nil {
		ErrResponse(ctx, http.StatusInternalServerError, ErrInternalServer.Error())
		return err
	}
	return ctx.Blob(http.StatusOK, "text/csv", file)
}
