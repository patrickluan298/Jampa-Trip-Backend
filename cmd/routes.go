package cmd

import (
	"github.com/jampa_trip/internal/app/handler"
	"github.com/jampa_trip/internal/app/middleware"
	"github.com/labstack/echo/v4"
)

func ConfigureRoutes(e *echo.Echo) {
	authHandler := handler.AuthHandler{}

	// Auth
	e.POST("/jampa-trip/api/v1/auth/login", middleware.ValidateJSONMiddleware(authHandler.Login))
	e.POST("/jampa-trip/api/v1/auth/register", middleware.ValidateJSONMiddleware(authHandler.Register))
	e.POST("/jampa-trip/api/v1/auth/logout", middleware.ValidateJSONMiddleware(authHandler.Logout))
	e.POST("/jampa-trip/api/v1/auth/refresh", middleware.ValidateJSONMiddleware(authHandler.Refresh))
}
