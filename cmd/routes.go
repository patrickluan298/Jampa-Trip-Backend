package cmd

import (
	"github.com/jampa_trip/internal/app/handler"
	"github.com/jampa_trip/internal/app/middleware"
	"github.com/labstack/echo/v4"
)

func ConfigureRoutes(e *echo.Echo) {

	// Health check
	e.GET("/health-check", handler.HealthCheckResponse{}.HealthCheck)
	e.HEAD("/health-check", handler.HealthCheckResponse{}.HealthCheck)

	// Empresas de turismo
	e.POST("/jampa-trip/api/v1/fornecedores/login", middleware.ValidateJSONMiddleware(handler.FornecedorHandler{}.Login))
	e.POST("/jampa-trip/api/v1/fornecedores/cadastrar", middleware.ValidateJSONMiddleware(handler.FornecedorHandler{}.Cadastrar))
	e.POST("/jampa-trip/api/v1/fornecedores/logout", middleware.ValidateJSONMiddleware(handler.FornecedorHandler{}.Logout))
	e.POST("/jampa-trip/api/v1/fornecedores/refresh", middleware.ValidateJSONMiddleware(handler.FornecedorHandler{}.Refresh))

	// Consumidores de serviços
	e.POST("/jampa-trip/api/v1/clientes/login", middleware.ValidateJSONMiddleware(handler.FornecedorHandler{}.Login))
	e.POST("/jampa-trip/api/v1/clientes/cadastrar", middleware.ValidateJSONMiddleware(handler.FornecedorHandler{}.Cadastrar))
	e.POST("/jampa-trip/api/v1/clientes/logout", middleware.ValidateJSONMiddleware(handler.FornecedorHandler{}.Logout))
	e.POST("/jampa-trip/api/v1/clientes/refresh", middleware.ValidateJSONMiddleware(handler.FornecedorHandler{}.Refresh))
}
