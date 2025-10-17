package main

import (
	"github.com/jampa_trip/internal/handler"
	"github.com/jampa_trip/pkg/middleware"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func ConfigureRoutes(e *echo.Echo) {

	// DOCUMENTATION
	e.GET("/docs/*", echoSwagger.WrapHandler)

	// HEALTH-CHECK
	e.GET("/health-check", handler.HealthCheckResponse{}.HealthCheck)
	e.HEAD("/health-check", handler.HealthCheckResponse{}.HealthCheck)

	// AUTHENTICATION
	e.POST("/jampa-trip/api/v1/login", handler.LoginHandler{}.Login)
	e.POST("/jampa-trip/api/v1/refresh", handler.RefreshHandler{}.RefreshToken)

	// PUBLIC REGISTER ROUTES
	e.POST("/jampa-trip/api/v1/companies", handler.CompanyHandler{}.Create)
	e.POST("/jampa-trip/api/v1/clients", handler.ClientHandler{}.Create)

	// PROTECTED GROUP – all routes below require JWT authentication
	protected := e.Group("/jampa-trip/api/v1")
	protected.Use(middleware.JWTMiddleware())

	// COMPANIES
	protected.PATCH("/companies/:id", handler.CompanyHandler{}.Update)
	protected.GET("/companies", handler.CompanyHandler{}.List)
	protected.GET("/companies/:id", handler.CompanyHandler{}.Get)

	// CLIENTS
	protected.PATCH("/clients/:id", handler.ClientHandler{}.Update)
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

	// IMAGE UPLOAD
	protected.POST("/upload/images", handler.ImageHandler{}.UploadImages)
	protected.GET("/upload/images", handler.ImageHandler{}.ListImages)
	protected.DELETE("/upload/images/:id", handler.ImageHandler{}.DeleteImage)
	protected.PUT("/upload/images/:id", handler.ImageHandler{}.UpdateImage)
	protected.POST("/upload/images/reorder", handler.ImageHandler{}.ReorderImages)
	protected.GET("/upload/images/:id/info", handler.ImageHandler{}.GetImageInfo)
	protected.POST("/upload/images/batch-delete", handler.ImageHandler{}.BatchDeleteImages)

	// FEEDBACK
	protected.POST("/feedback", handler.FeedbackHandler{}.Create)
	protected.GET("/feedback/:id", handler.FeedbackHandler{}.Get)
	protected.GET("/feedback", handler.FeedbackHandler{}.List)
	protected.PUT("/feedback/:id", handler.FeedbackHandler{}.Update)
	protected.GET("/feedback/average-rating", handler.FeedbackHandler{}.GetAverageRating)
	protected.GET("/feedback/rating-distribution", handler.FeedbackHandler{}.GetRatingDistribution)
	protected.GET("/feedback/recent", handler.FeedbackHandler{}.GetRecent)

	// RESERVATIONS
	// protected.POST("/reservations", handler.ReservaHandler{}.Create)
	// protected.GET("/reservations/:id", handler.ReservaHandler{}.Get)
	// protected.GET("/reservations", handler.ReservaHandler{}.List)
	// protected.PUT("/reservations/:id", handler.ReservaHandler{}.Update)
	// protected.PUT("/reservations/:id/cancel", handler.ReservaHandler{}.Cancel)
	// protected.GET("/reservations/upcoming", handler.ReservaHandler{}.GetUpcoming)
	// protected.GET("/reservations/history", handler.ReservaHandler{}.GetHistory)
}
