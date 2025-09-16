package cmd

import (
	"github.com/jampa_trip/internal/app/handler"
	"github.com/labstack/echo/v4"
)

func ConfigureRoutes(e *echo.Echo) {

	// Health check
	e.GET("/health-check", handler.HealthCheckResponse{}.HealthCheck)
	e.HEAD("/health-check", handler.HealthCheckResponse{}.HealthCheck)

	// Empresas de turismo
	e.POST("/jampa-trip/api/v1/fornecedores/login", handler.FornecedorHandler{}.Login)
	e.POST("/jampa-trip/api/v1/fornecedores/cadastrar", handler.FornecedorHandler{}.Cadastrar)
	e.POST("/jampa-trip/api/v1/fornecedores/logout", handler.FornecedorHandler{}.Logout)
	e.POST("/jampa-trip/api/v1/fornecedores/refresh", handler.FornecedorHandler{}.Refresh)

	// Consumidores de serviços
	e.POST("/jampa-trip/api/v1/clientes/login", handler.FornecedorHandler{}.Login)
	e.POST("/jampa-trip/api/v1/clientes/cadastrar", handler.FornecedorHandler{}.Cadastrar)
	e.POST("/jampa-trip/api/v1/clientes/logout", handler.FornecedorHandler{}.Logout)
	e.POST("/jampa-trip/api/v1/clientes/refresh", handler.FornecedorHandler{}.Refresh)
}
