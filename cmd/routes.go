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

	// Pagamentos
	e.POST("/jampa-trip/api/v1/pagamentos", handler.PagamentoHandler{}.Create)
	e.GET("/jampa-trip/api/v1/pagamentos/:id", handler.PagamentoHandler{}.Get)
	e.GET("/jampa-trip/api/v1/pagamentos", handler.PagamentoHandler{}.List)
	e.PUT("/jampa-trip/api/v1/pagamentos/:id", handler.PagamentoHandler{}.Update)
	e.DELETE("/jampa-trip/api/v1/pagamentos/:id", handler.PagamentoHandler{}.Delete)

	// Reservas
	e.POST("/jampa-trip/api/v1/reservas", handler.ReservaHandler{}.Create)
	e.GET("/jampa-trip/api/v1/reservas/:id", handler.ReservaHandler{}.Get)
	e.GET("/jampa-trip/api/v1/reservas", handler.ReservaHandler{}.List)
	e.PUT("/jampa-trip/api/v1/reservas/:id", handler.ReservaHandler{}.Update)
	e.PUT("/jampa-trip/api/v1/reservas/:id/cancelar", handler.ReservaHandler{}.Cancel)
	e.GET("/jampa-trip/api/v1/reservas/futuras", handler.ReservaHandler{}.GetUpcoming)
	e.GET("/jampa-trip/api/v1/reservas/historico", handler.ReservaHandler{}.GetHistory)

	// Feedback e comentários
	e.POST("/jampa-trip/api/v1/feedback", handler.FeedbackHandler{}.Create)
	e.GET("/jampa-trip/api/v1/feedback/:id", handler.FeedbackHandler{}.Get)
	e.GET("/jampa-trip/api/v1/feedback", handler.FeedbackHandler{}.List)
	e.PUT("/jampa-trip/api/v1/feedback/:id", handler.FeedbackHandler{}.Update)
	e.GET("/jampa-trip/api/v1/feedback/avaliacao-media", handler.FeedbackHandler{}.GetAverageRating)
	e.GET("/jampa-trip/api/v1/feedback/distribuicao-notas", handler.FeedbackHandler{}.GetRatingDistribution)
	e.GET("/jampa-trip/api/v1/feedback/recentes", handler.FeedbackHandler{}.GetRecent)
}
