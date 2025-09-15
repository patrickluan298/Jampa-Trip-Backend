package handler

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

// HealthCheckResponse referente ao response da função
type HealthCheckResponse struct {
	Status             int    `json:"status"`
	Message            string `json:"message"`
	NameApplication    string `json:"nameApplication"`
	VersionApplication string `json:"versionApplication"`
}

func (h HealthCheckResponse) HealthCheck(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, &HealthCheckResponse{
		Status:             http.StatusOK,
		Message:            "Aplicação up e em pleno funcionamento",
		NameApplication:    "Jampa-Trip",
		VersionApplication: os.Getenv("VERSION_APPLICATION"),
	})
}
