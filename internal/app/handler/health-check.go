package handler

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

// HealthCheckResponse referente ao response da função
type HealthCheckResponse struct {
	Nome     string `json:"nome"`
	Versao   string `json:"versao"`
	Mensagem string `json:"mensagem"`
}

func (h HealthCheckResponse) HealthCheck(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, &HealthCheckResponse{
		Nome:     "Jampa-Trip",
		Versao:   os.Getenv("VERSION_APPLICATION"),
		Mensagem: "Aplicação up e em pleno funcionamento",
	})
}
