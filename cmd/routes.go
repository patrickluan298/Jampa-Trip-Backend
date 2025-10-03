package main

import (
	"github.com/jampa_trip/internal/handler"
	"github.com/labstack/echo/v4"
)

func ConfigureRoutes(e *echo.Echo) {

	// HEALTH-CHECK
	e.GET("/health-check", handler.HealthCheckResponse{}.HealthCheck)
	e.HEAD("/health-check", handler.HealthCheckResponse{}.HealthCheck)

	// EMPRESAS
	e.POST("/jampa-trip/api/v1/empresas/login", handler.EmpresaHandler{}.Login)
	e.POST("/jampa-trip/api/v1/empresas/cadastrar", handler.EmpresaHandler{}.Create)
	e.PUT("/jampa-trip/api/v1/empresas/atualizar/:id", handler.EmpresaHandler{}.Update)
	e.GET("/jampa-trip/api/v1/empresas/listar", handler.EmpresaHandler{}.List)
	e.GET("/jampa-trip/api/v1/empresas/obter/:id", handler.EmpresaHandler{}.Get)

	// CLIENTES
	e.POST("/jampa-trip/api/v1/clientes/login", handler.ClienteHandler{}.Login)
	e.POST("/jampa-trip/api/v1/clientes/cadastrar", handler.ClienteHandler{}.Create)
	e.PUT("/jampa-trip/api/v1/clientes/atualizar/:id", handler.ClienteHandler{}.Update)
	e.GET("/jampa-trip/api/v1/clientes/listar", handler.ClienteHandler{}.List)
	e.GET("/jampa-trip/api/v1/clientes/obter/:id", handler.ClienteHandler{}.Get)

	// CARTÕES
	e.POST("/jampa-trip/api/v1/clientes/:customer_id/cartoes", handler.CartaoHandler{}.Create)
	e.GET("/jampa-trip/api/v1/clientes/:customer_id/cartoes", handler.CartaoHandler{}.List)
	e.GET("/jampa-trip/api/v1/clientes/:customer_id/cartoes/:card_id", handler.CartaoHandler{}.Get)
	e.PUT("/jampa-trip/api/v1/clientes/:customer_id/cartoes/:card_id", handler.CartaoHandler{}.Update)
	e.DELETE("/jampa-trip/api/v1/clientes/:customer_id/cartoes/:card_id", handler.CartaoHandler{}.Delete)

	// MÉTODOS DE PAGAMENTO
	e.POST("/jampa-trip/api/v1/pagamentos/cartao-credito", handler.PagamentoHandler{}.CriarPagamentoCartaoCredito)
	e.POST("/jampa-trip/api/v1/pagamentos/cartao-debito", handler.PagamentoHandler{}.CriarPagamentoCartaoDebito)
	e.POST("/jampa-trip/api/v1/pagamentos/pix", handler.PagamentoHandler{}.CriarPagamentoPIX)
	e.GET("/jampa-trip/api/v1/pagamentos", handler.PagamentoHandler{}.BuscarPagamentos)
	e.GET("/jampa-trip/api/v1/pagamentos/:id", handler.PagamentoHandler{}.ObterPagamento)
	e.PUT("/jampa-trip/api/v1/pagamentos/:id", handler.PagamentoHandler{}.AtualizarPagamento)

	// RESERVAS
	// e.POST("/jampa-trip/api/v1/reservas", handler.ReservaHandler{}.Create)
	// e.GET("/jampa-trip/api/v1/reservas/:id", handler.ReservaHandler{}.Get)
	// e.GET("/jampa-trip/api/v1/reservas", handler.ReservaHandler{}.List)
	// e.PUT("/jampa-trip/api/v1/reservas/:id", handler.ReservaHandler{}.Update)
	// e.PUT("/jampa-trip/api/v1/reservas/:id/cancelar", handler.ReservaHandler{}.Cancel)
	// e.GET("/jampa-trip/api/v1/reservas/futuras", handler.ReservaHandler{}.GetUpcoming)
	// e.GET("/jampa-trip/api/v1/reservas/historico", handler.ReservaHandler{}.GetHistory)

	// FEEDBACK E COMENTÁRIOS
	// e.POST("/jampa-trip/api/v1/feedback", handler.FeedbackHandler{}.Create)
	// e.GET("/jampa-trip/api/v1/feedback/:id", handler.FeedbackHandler{}.Get)
	// e.GET("/jampa-trip/api/v1/feedback", handler.FeedbackHandler{}.List)
	// e.PUT("/jampa-trip/api/v1/feedback/:id", handler.FeedbackHandler{}.Update)
	// e.GET("/jampa-trip/api/v1/feedback/avaliacao-media", handler.FeedbackHandler{}.GetAverageRating)
	// e.GET("/jampa-trip/api/v1/feedback/distribuicao-notas", handler.FeedbackHandler{}.GetRatingDistribution)
	// e.GET("/jampa-trip/api/v1/feedback/recentes", handler.FeedbackHandler{}.GetRecent)
}
