package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/jampa_trip/pkg/database"
	"github.com/jampa_trip/pkg/util"
	"github.com/redis/go-redis/v9"
)

// TokenStore - interface para gerenciamento de tokens
type TokenStore interface {
	StoreAccessToken(userID int, userType, token string) error
	StoreRefreshToken(userID int, userType, token string) error
	ValidateAccessToken(userID int, userType, token string) error
	ValidateRefreshToken(userID int, userType, token string) error
	DeleteTokens(userID int, userType string) error
}

// RedisTokenStore - implementação do TokenStore usando Redis
type RedisTokenStore struct {
	client *redis.Client
}

// NewRedisTokenStore - cria uma nova instância do RedisTokenStore
func NewRedisTokenStore() *RedisTokenStore {
	return &RedisTokenStore{
		client: database.RedisClient,
	}
}

// StoreAccessToken - armazena access token no Redis com TTL de 15 minutos
func (r *RedisTokenStore) StoreAccessToken(userID int, userType, token string) error {
	ctx := context.Background()
	key := fmt.Sprintf("access_token:%d:%s", userID, userType)

	duration, err := time.ParseDuration(database.Config.JWTAccessTokenExpiration)
	if err != nil {
		return util.WrapError("erro ao parsear duração do access token", err, 500)
	}

	err = r.client.Set(ctx, key, token, duration).Err()
	if err != nil {
		return util.WrapError("erro ao armazenar access token no Redis", err, 500)
	}

	return nil
}

// StoreRefreshToken - armazena refresh token no Redis com TTL de 7 dias
func (r *RedisTokenStore) StoreRefreshToken(userID int, userType, token string) error {
	ctx := context.Background()
	key := fmt.Sprintf("refresh_token:%d:%s", userID, userType)

	duration, err := time.ParseDuration(database.Config.JWTRefreshTokenExpiration)
	if err != nil {
		return util.WrapError("erro ao parsear duração do refresh token", err, 500)
	}

	err = r.client.Set(ctx, key, token, duration).Err()
	if err != nil {
		return util.WrapError("erro ao armazenar refresh token no Redis", err, 500)
	}

	return nil
}

// ValidateAccessToken - verifica se o access token existe no Redis
func (r *RedisTokenStore) ValidateAccessToken(userID int, userType, token string) error {
	ctx := context.Background()
	key := fmt.Sprintf("access_token:%d:%s", userID, userType)

	storedToken, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return util.WrapError("access token não encontrado no Redis", err, 401)
	}

	if storedToken != token {
		return util.WrapError("access token inválido", nil, 401)
	}

	return nil
}

// ValidateRefreshToken - verifica se o refresh token existe no Redis
func (r *RedisTokenStore) ValidateRefreshToken(userID int, userType, token string) error {
	ctx := context.Background()
	key := fmt.Sprintf("refresh_token:%d:%s", userID, userType)

	storedToken, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return util.WrapError("refresh token não encontrado no Redis", err, 401)
	}

	if storedToken != token {
		return util.WrapError("refresh token inválido", nil, 401)
	}

	return nil
}

// DeleteTokens - remove tokens do Redis (para logout)
func (r *RedisTokenStore) DeleteTokens(userID int, userType string) error {
	ctx := context.Background()
	accessKey := fmt.Sprintf("access_token:%d:%s", userID, userType)
	refreshKey := fmt.Sprintf("refresh_token:%d:%s", userID, userType)

	err := r.client.Del(ctx, accessKey).Err()
	if err != nil {
		return util.WrapError("erro ao remover access token do Redis", err, 500)
	}

	err = r.client.Del(ctx, refreshKey).Err()
	if err != nil {
		return util.WrapError("erro ao remover refresh token do Redis", err, 500)
	}

	return nil
}

// StoreTokenPair - armazena ambos os tokens (access e refresh)
func (r *RedisTokenStore) StoreTokenPair(userID int, userType string, accessToken, refreshToken string) error {
	err := r.StoreAccessToken(userID, userType, accessToken)
	if err != nil {
		return err
	}

	err = r.StoreRefreshToken(userID, userType, refreshToken)
	if err != nil {
		r.DeleteTokens(userID, userType)
		return err
	}

	return nil
}
