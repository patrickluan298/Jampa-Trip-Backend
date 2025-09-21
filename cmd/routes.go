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
	e.POST("/jampa-trip/api/v1/fornecedores/atualizar", handler.FornecedorHandler{}.Atualizar)
	e.GET("/jampa-trip/api/v1/fornecedores/listar", handler.FornecedorHandler{}.Listar)

	// Consumidores de serviços
	e.POST("/jampa-trip/api/v1/clientes/login", handler.ClienteHandler{}.Login)
	e.POST("/jampa-trip/api/v1/clientes/cadastrar", handler.ClienteHandler{}.Cadastrar)
	e.POST("/jampa-trip/api/v1/clientes/atualizar", handler.ClienteHandler{}.Atualizar)
	e.GET("/jampa-trip/api/v1/clientes/listar", handler.ClienteHandler{}.Listar)
}
