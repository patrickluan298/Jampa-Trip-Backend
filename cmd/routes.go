package cmd

import (
	"github.com/jampa_trip/internal/app/handler"
	"github.com/labstack/echo/v4"
)

func ConfigureRoutes(e *echo.Echo) {

	// Health check
	e.GET("/health-check", handler.HealthCheckResponse{}.HealthCheck)
	e.HEAD("/health-check", handler.HealthCheckResponse{}.HealthCheck)

	// Empresas
	e.POST("/jampa-trip/api/v1/empresas/login", handler.EmpresaHandler{}.Login)
	e.POST("/jampa-trip/api/v1/empresas/cadastrar", handler.EmpresaHandler{}.Create)
	e.PUT("/jampa-trip/api/v1/empresas/atualizar/:id", handler.EmpresaHandler{}.Update)
	e.GET("/jampa-trip/api/v1/empresas/listar", handler.EmpresaHandler{}.List)
	e.GET("/jampa-trip/api/v1/empresas/obter/:id", handler.EmpresaHandler{}.Get)

	// Clientes
	e.POST("/jampa-trip/api/v1/clientes/login", handler.ClienteHandler{}.Login)
	e.POST("/jampa-trip/api/v1/clientes/cadastrar", handler.ClienteHandler{}.Create)
	e.PUT("/jampa-trip/api/v1/clientes/atualizar/:id", handler.ClienteHandler{}.Update)
	e.GET("/jampa-trip/api/v1/clientes/listar", handler.ClienteHandler{}.List)
	e.GET("/jampa-trip/api/v1/clientes/obter/:id", handler.ClienteHandler{}.Get)

	// Métodos de pagamento
	e.POST("/jampa-trip/api/v1/pagamentos", handler.ClienteHandler{}.Create)
	e.PUT("/jampa-trip/api/v1/pagamentos/:id", handler.ClienteHandler{}.Update)
	e.GET("/jampa-trip/api/v1/pagamentos", handler.ClienteHandler{}.List)
	e.DELETE("/jampa-trip/api/v1/pagamentos/:id", handler.ClienteHandler{}.List)

	// Reservas
	e.GET("/jampa-trip/api/v1/reservas", handler.ClienteHandler{}.List)
	e.GET("/jampa-trip/api/v1/reservas/:id", handler.ClienteHandler{}.Get)
	e.PUT("/jampa-trip/api/v1/reservas/:id/cancelar", handler.ClienteHandler{}.Update)

	// Histórico de passeios
	e.GET("/jampa-trip/api/v1/historico", handler.ClienteHandler{}.List)

	// Feedback e comentários
	e.GET("/jampa-trip/api/v1/feedback", handler.ClienteHandler{}.List)
	e.POST("/jampa-trip/api/v1/feedback", handler.ClienteHandler{}.Create)
}
