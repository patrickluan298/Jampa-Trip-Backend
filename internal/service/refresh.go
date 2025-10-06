package service

import (
	"net/http"

	"github.com/jampa_trip/internal/contract"
	"github.com/jampa_trip/pkg/auth"
	"github.com/jampa_trip/pkg/util"
)

// RefreshService - objeto de contexto para refresh token
type RefreshService struct{}

// RefreshServiceNew - construtor do objeto
func RefreshServiceNew() *RefreshService {
	return &RefreshService{}
}

// RefreshToken - renova o par de tokens
func (receiver *RefreshService) RefreshToken(request *contract.RefreshTokenRequest) (*contract.RefreshTokenResponse, error) {
	claims, err := auth.ValidateToken(request.RefreshToken)
	if err != nil {
		return nil, util.WrapError("refresh token inválido", err, http.StatusUnauthorized)
	}

	if auth.IsTokenExpired(claims) {
		return nil, util.WrapError("refresh token expirado", nil, http.StatusUnauthorized)
	}

	tokenStore := auth.NewRedisTokenStore()

	err = tokenStore.ValidateRefreshToken(claims.UserID, claims.UserType, request.RefreshToken)
	if err != nil {
		return nil, util.WrapError("refresh token não encontrado ou inválido", err, http.StatusUnauthorized)
	}

	newTokenPair, err := auth.GenerateTokenPair(claims.UserID, claims.UserType, claims.Email)
	if err != nil {
		return nil, util.WrapError("erro ao gerar novos tokens JWT", err, http.StatusInternalServerError)
	}

	err = tokenStore.DeleteTokens(claims.UserID, claims.UserType)
	if err != nil {
		return nil, util.WrapError("erro ao remover tokens antigos do Redis", err, http.StatusInternalServerError)
	}

	err = tokenStore.StoreTokenPair(claims.UserID, claims.UserType, newTokenPair.AccessToken, newTokenPair.RefreshToken)
	if err != nil {
		return nil, util.WrapError("erro ao armazenar novos tokens no Redis", err, http.StatusInternalServerError)
	}

	response := &contract.RefreshTokenResponse{
		AccessToken:  newTokenPair.AccessToken,
		RefreshToken: newTokenPair.RefreshToken,
		ExpiresIn:    newTokenPair.ExpiresIn,
	}

	return response, nil
}
