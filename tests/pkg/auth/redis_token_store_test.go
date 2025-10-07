package auth

import (
	"testing"

	"github.com/jampa_trip/pkg/auth"
)

func TestNewRedisTokenStore(t *testing.T) {
	// Skip this test as it requires global database.RedisClient
	t.Skip("Skipping due to global dependency - requires dependency injection")
	
	store := auth.NewRedisTokenStore()
	if store == nil {
		t.Errorf("NewRedisTokenStore() returned nil")
	}
}

func TestRedisTokenStore_StoreAccessToken(t *testing.T) {
	// Skip this test as it requires global database.RedisClient
	t.Skip("Skipping due to global dependency - requires dependency injection")
}

func TestRedisTokenStore_StoreRefreshToken(t *testing.T) {
	// Skip this test as it requires global database.RedisClient
	t.Skip("Skipping due to global dependency - requires dependency injection")
}

func TestRedisTokenStore_ValidateAccessToken(t *testing.T) {
	// Skip this test as it requires global database.RedisClient
	t.Skip("Skipping due to global dependency - requires dependency injection")
}

func TestRedisTokenStore_ValidateRefreshToken(t *testing.T) {
	// Skip this test as it requires global database.RedisClient
	t.Skip("Skipping due to global dependency - requires dependency injection")
}

func TestRedisTokenStore_DeleteTokens(t *testing.T) {
	// Skip this test as it requires global database.RedisClient
	t.Skip("Skipping due to global dependency - requires dependency injection")
}