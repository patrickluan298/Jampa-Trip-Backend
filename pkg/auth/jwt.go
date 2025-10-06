package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jampa_trip/pkg/database"
	"github.com/jampa_trip/pkg/util"
)

// JWTClaims - estrutura de claims do JWT
type JWTClaims struct {
	UserID   int    `json:"user_id"`
	UserType string `json:"user_type"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

// TokenPair - par de tokens (access e refresh)
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

// GenerateTokenPair - gera um par de tokens (access e refresh)
func GenerateTokenPair(userID int, userType, email string) (*TokenPair, error) {

	accessDuration, err := time.ParseDuration(database.Config.JWTAccessTokenExpiration)
	if err != nil {
		return nil, util.WrapError("erro ao parsear duração do access token", err, 500)
	}

	refreshDuration, err := time.ParseDuration(database.Config.JWTRefreshTokenExpiration)
	if err != nil {
		return nil, util.WrapError("erro ao parsear duração do refresh token", err, 500)
	}

	now := time.Now()
	accessExpiresAt := now.Add(accessDuration)
	refreshExpiresAt := now.Add(refreshDuration)

	accessClaims := JWTClaims{
		UserID:   userID,
		UserType: userType,
		Email:    email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessExpiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "jampa-trip",
			Subject:   fmt.Sprintf("%d", userID),
		},
	}

	refreshClaims := JWTClaims{
		UserID:   userID,
		UserType: userType,
		Email:    email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshExpiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "jampa-trip",
			Subject:   fmt.Sprintf("%d", userID),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(database.Config.JWTSecret))
	if err != nil {
		return nil, util.WrapError("erro ao assinar access token", err, 500)
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(database.Config.JWTSecret))
	if err != nil {
		return nil, util.WrapError("erro ao assinar refresh token", err, 500)
	}

	return &TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		ExpiresIn:    accessExpiresAt.Unix(),
	}, nil
}

// ValidateToken - valida e extrai claims do token
func ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := ParseToken(tokenString)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return nil, util.WrapError("erro ao extrair claims do token", nil, 401)
	}

	return claims, nil
}

// ParseToken - faz o parse do token JWT
func ParseToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, util.WrapError("método de assinatura inesperado", nil, 401)
		}
		return []byte(database.Config.JWTSecret), nil
	})

	if err != nil {
		return nil, util.WrapError("erro ao fazer parse do token", err, 401)
	}

	if !token.Valid {
		return nil, util.WrapError("token inválido", nil, 401)
	}

	return token, nil
}

// IsTokenExpired - verifica se o token está expirado
func IsTokenExpired(claims *JWTClaims) bool {
	return claims.ExpiresAt.Before(time.Now())
}
