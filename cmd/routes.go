package main

import (
	"github.com/jampa_trip/internal/handler"
	"github.com/labstack/echo/v4"
)

func ConfigureRoutes(e *echo.Echo) {

	// HEALTH-CHECK
	e.GET("/health-check", handler.HealthCheckResponse{}.HealthCheck)
	e.HEAD("/health-check", handler.HealthCheckResponse{}.HealthCheck)

	// LOGIN
	e.POST("/jampa-trip/api/v1/login", handler.LoginHandler{}.Login)

	// COMPANIES
	e.POST("/jampa-trip/api/v1/companies", handler.CompanyHandler{}.Create)
	e.PUT("/jampa-trip/api/v1/companies/:id", handler.CompanyHandler{}.Update)
	e.GET("/jampa-trip/api/v1/companies", handler.CompanyHandler{}.List)
	e.GET("/jampa-trip/api/v1/companies/:id", handler.CompanyHandler{}.Get)

	// CLIENTS
	e.POST("/jampa-trip/api/v1/clients", handler.ClientHandler{}.Create)
	e.PUT("/jampa-trip/api/v1/clients/:id", handler.ClientHandler{}.Update)
	e.GET("/jampa-trip/api/v1/clients", handler.ClientHandler{}.List)
	e.GET("/jampa-trip/api/v1/clients/:id", handler.ClientHandler{}.Get)

	// CARDS
	e.POST("/jampa-trip/api/v1/clients/:customer_id/cards", handler.CartaoHandler{}.Create)
	e.GET("/jampa-trip/api/v1/clients/:customer_id/cards", handler.CartaoHandler{}.List)
	e.GET("/jampa-trip/api/v1/clients/:customer_id/cards/:card_id", handler.CartaoHandler{}.Get)
	e.PUT("/jampa-trip/api/v1/clients/:customer_id/cards/:card_id", handler.CartaoHandler{}.Update)
	e.DELETE("/jampa-trip/api/v1/clients/:customer_id/cards/:card_id", handler.CartaoHandler{}.Delete)

	// PAYMENT METHODS
	e.POST("/jampa-trip/api/v1/payments/credit-card", handler.PagamentoHandler{}.CriarPagamentoCartaoCredito)
	e.POST("/jampa-trip/api/v1/payments/debit-card", handler.PagamentoHandler{}.CriarPagamentoCartaoDebito)
	e.POST("/jampa-trip/api/v1/payments/pix", handler.PagamentoHandler{}.CriarPagamentoPIX)
	e.GET("/jampa-trip/api/v1/payments", handler.PagamentoHandler{}.BuscarPagamentos)
	e.GET("/jampa-trip/api/v1/payments/:id", handler.PagamentoHandler{}.ObterPagamento)
	e.PUT("/jampa-trip/api/v1/payments/:id", handler.PagamentoHandler{}.AtualizarPagamento)

	// RESERVATIONS
	// e.POST("/jampa-trip/api/v1/reservations", handler.ReservaHandler{}.Create)
	// e.GET("/jampa-trip/api/v1/reservations/:id", handler.ReservaHandler{}.Get)
	// e.GET("/jampa-trip/api/v1/reservations", handler.ReservaHandler{}.List)
	// e.PUT("/jampa-trip/api/v1/reservations/:id", handler.ReservaHandler{}.Update)
	// e.PUT("/jampa-trip/api/v1/reservations/:id/cancel", handler.ReservaHandler{}.Cancel)
	// e.GET("/jampa-trip/api/v1/reservations/upcoming", handler.ReservaHandler{}.GetUpcoming)
	// e.GET("/jampa-trip/api/v1/reservations/history", handler.ReservaHandler{}.GetHistory)

	// FEEDBACK AND COMMENTS
	// e.POST("/jampa-trip/api/v1/feedback", handler.FeedbackHandler{}.Create)
	// e.GET("/jampa-trip/api/v1/feedback/:id", handler.FeedbackHandler{}.Get)
	// e.GET("/jampa-trip/api/v1/feedback", handler.FeedbackHandler{}.List)
	// e.PUT("/jampa-trip/api/v1/feedback/:id", handler.FeedbackHandler{}.Update)
	// e.GET("/jampa-trip/api/v1/feedback/average-rating", handler.FeedbackHandler{}.GetAverageRating)
	// e.GET("/jampa-trip/api/v1/feedback/rating-distribution", handler.FeedbackHandler{}.GetRatingDistribution)
	// e.GET("/jampa-trip/api/v1/feedback/recent", handler.FeedbackHandler{}.GetRecent)
}
