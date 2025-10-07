package main

import (
	"github.com/jampa_trip/internal/handler"
	"github.com/jampa_trip/pkg/middleware"
	"github.com/labstack/echo/v4"
)

func ConfigureRoutes(e *echo.Echo) {

	// HEALTH-CHECK
	e.GET("/health-check", handler.HealthCheckResponse{}.HealthCheck)
	e.HEAD("/health-check", handler.HealthCheckResponse{}.HealthCheck)

	// AUTHENTICATION
	e.POST("/jampa-trip/api/v1/login", handler.LoginHandler{}.Login)
	e.POST("/jampa-trip/api/v1/refresh", handler.RefreshHandler{}.RefreshToken)

	// GRUPO PROTEGIDO - todas as rotas abaixo precisam de autenticação JWT
	protected := e.Group("/jampa-trip/api/v1")
	protected.Use(middleware.JWTMiddleware())

	// COMPANIES
	protected.POST("/companies", handler.CompanyHandler{}.Create)
	protected.PUT("/companies/:id", handler.CompanyHandler{}.Update)
	protected.GET("/companies", handler.CompanyHandler{}.List)
	protected.GET("/companies/:id", handler.CompanyHandler{}.Get)

	// CLIENTS
	protected.POST("/clients", handler.ClientHandler{}.Create)
	protected.PUT("/clients/:id", handler.ClientHandler{}.Update)
	protected.GET("/clients", handler.ClientHandler{}.List)
	protected.GET("/clients/:id", handler.ClientHandler{}.Get)

	// CARDS
	protected.POST("/clients/:customer_id/cards", handler.CardHandler{}.Create)
	protected.GET("/clients/:customer_id/cards", handler.CardHandler{}.List)
	protected.GET("/clients/:customer_id/cards/:card_id", handler.CardHandler{}.Get)
	protected.PUT("/clients/:customer_id/cards/:card_id", handler.CardHandler{}.Update)
	protected.DELETE("/clients/:customer_id/cards/:card_id", handler.CardHandler{}.Delete)

	// PAYMENT METHODS
	protected.POST("/payments/credit-card", handler.PaymentHandler{}.CreateCreditCardPayment)
	protected.POST("/payments/debit-card", handler.PaymentHandler{}.CreateDebitCardPayment)
	protected.POST("/payments/pix", handler.PaymentHandler{}.CreatePIXPayment)
	protected.GET("/payments", handler.PaymentHandler{}.List)
	protected.GET("/payments/:id", handler.PaymentHandler{}.Get)
	protected.PUT("/payments/:id", handler.PaymentHandler{}.Update)

	// TOURS
	protected.POST("/tours", handler.TourHandler{}.Create)
	protected.GET("/tours", handler.TourHandler{}.List)
	protected.PUT("/tours/:id", handler.TourHandler{}.Update)
	protected.DELETE("/tours/:id", handler.TourHandler{}.Delete)
	protected.GET("/tours/my-tours", handler.TourHandler{}.GetMyTours)

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
