package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// SetupMiddlewares - configura os middlewares
func SetupMiddlewares(e *echo.Echo) {
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.Gzip())
}

// ValidateJSONMiddleware - verifica se o JSON da requisição é válido
func ValidateJSONMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var data map[string]interface{}
		if err := c.Bind(&data); err != nil {
			return c.JSON(400, map[string]string{"error": "JSON inválido"})
		}

		return next(c)
	}
}
