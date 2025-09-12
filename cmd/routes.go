package cmd

import (
	"github.com/jampa_trip/internal/app/handler"
	"github.com/jampa_trip/internal/app/middleware"
	"github.com/labstack/echo/v4"
)

func ConfigureRoutes(e *echo.Echo) {

	// Empresas de turismo
	e.POST("/jampa-trip/api/v1/fornecedores/login", middleware.ValidateJSONMiddleware(handler.FornecedorHandler{}.Login))
	e.POST("/jampa-trip/api/v1/fornecedores/register", middleware.ValidateJSONMiddleware(handler.FornecedorHandler{}.Register))
	e.POST("/jampa-trip/api/v1/fornecedores/logout", middleware.ValidateJSONMiddleware(handler.FornecedorHandler{}.Logout))
	e.POST("/jampa-trip/api/v1/fornecedores/refresh", middleware.ValidateJSONMiddleware(handler.FornecedorHandler{}.Refresh))

	// Cconsumidores de serviços
	e.POST("/jampa-trip/api/v1/clientes/login", middleware.ValidateJSONMiddleware(handler.FornecedorHandler{}.Login))
	e.POST("/jampa-trip/api/v1/clientes/register", middleware.ValidateJSONMiddleware(handler.FornecedorHandler{}.Register))
	e.POST("/jampa-trip/api/v1/clientes/logout", middleware.ValidateJSONMiddleware(handler.FornecedorHandler{}.Logout))
	e.POST("/jampa-trip/api/v1/clientes/refresh", middleware.ValidateJSONMiddleware(handler.FornecedorHandler{}.Refresh))
}
