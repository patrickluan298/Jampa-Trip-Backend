package middleware

import (
	"net/http"
	"strings"

	"github.com/jampa_trip/pkg/auth"
	"github.com/jampa_trip/pkg/util"
	"github.com/jampa_trip/pkg/webserver"
	"github.com/labstack/echo/v4"
)

// JWTMiddleware - middleware para validação de JWT
func JWTMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return webserver.ErrorResponse(c, util.WrapError("token de autorização não fornecido", nil, http.StatusUnauthorized))
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return webserver.ErrorResponse(c, util.WrapError("formato de token inválido. Use: Bearer <token>", nil, http.StatusUnauthorized))
			}

			tokenString := parts[1]

			claims, err := auth.ValidateToken(tokenString)
			if err != nil {
				return webserver.ErrorResponse(c, err)
			}

			if auth.IsTokenExpired(claims) {
				return webserver.ErrorResponse(c, util.WrapError("token expirado", nil, http.StatusUnauthorized))
			}

			tokenStore := auth.NewRedisTokenStore()

			err = tokenStore.ValidateAccessToken(claims.UserID, claims.UserType, tokenString)
			if err != nil {
				return webserver.ErrorResponse(c, util.WrapError("token não encontrado ou inválido", err, http.StatusUnauthorized))
			}

			c.Set("user_id", claims.UserID)
			c.Set("user_type", claims.UserType)
			c.Set("user_email", claims.Email)
			c.Set("jwt_claims", claims)

			return next(c)
		}
	}
}

// GetUserID - extrai o ID do usuário do contexto
func GetUserID(c echo.Context) int {
	userID, ok := c.Get("user_id").(int)
	if !ok {
		return 0
	}
	return userID
}

// GetUserType - extrai o tipo do usuário do contexto
func GetUserType(c echo.Context) string {
	userType, ok := c.Get("user_type").(string)
	if !ok {
		return ""
	}
	return userType
}

// GetUserEmail - extrai o email do usuário do contexto
func GetUserEmail(c echo.Context) string {
	userEmail, ok := c.Get("user_email").(string)
	if !ok {
		return ""
	}
	return userEmail
}

// GetJWTClaims - extrai as claims completas do contexto
func GetJWTClaims(c echo.Context) *auth.JWTClaims {
	claims, ok := c.Get("jwt_claims").(*auth.JWTClaims)
	if !ok {
		return nil
	}
	return claims
}
