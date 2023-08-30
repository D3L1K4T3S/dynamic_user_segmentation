package v1

import (
	"dynamic-user-segmentation/internal/service"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/swaggo/echo-swagger"
	"net/http"
)

func NewRouter(handler *echo.Echo, services *service.Services) {
	handler.Use(middleware.Recover())
	handler.GET("/swagger/*", echoSwagger.WrapHandler)
	handler.GET("/health", func(ctx echo.Context) error {
		return ctx.NoContent(http.StatusOK)
	})

	auth := handler.Group("/auth")
	newAuthRoutes(auth, services.Auth)

	authMiddleware := &AuthMiddleware{services.Auth}
	v1 := handler.Group("/api/v1", authMiddleware.CheckUser)
	newActionsRoutes(v1.Group("/actions"), services.Actions)
	newConsumersRoutes(v1.Group("/consumers"), services.Consumers)
	newOperationsRoutes(v1.Group("/operations"), services.Operations)
	newSegmentRoutes(v1.Group("/segments"), services.Segments)
}
