package testutils

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// MockFactory provides a unified way to create mocks for testing
type MockFactory struct {
	DB        *gorm.DB
	Redis     *redis.Client
	SQLMock   sqlmock.Sqlmock
	MiniRedis *miniredis.Miniredis
	T         *testing.T
}

// NewMockFactory creates a new MockFactory with all necessary mocks
func NewMockFactory(t *testing.T) *MockFactory {
	t.Helper()

	factory := &MockFactory{
		T: t,
	}

	// Setup database mock
	db, sqlMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to create GORM DB: %v", err)
	}

	factory.DB = gormDB
	factory.SQLMock = sqlMock

	// Setup Redis mock
	mr := miniredis.RunT(t)
	redisClient := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	factory.Redis = redisClient
	factory.MiniRedis = mr

	return factory
}

// Cleanup cleans up all mocks
func (f *MockFactory) Cleanup() {
	if f.MiniRedis != nil {
		f.MiniRedis.Close()
	}
	if f.Redis != nil {
		f.Redis.Close()
	}
}

// ExpectationsWereMet checks if all SQL mock expectations were met
func (f *MockFactory) ExpectationsWereMet() {
	if err := f.SQLMock.ExpectationsWereMet(); err != nil {
		f.T.Errorf("SQL mock expectations were not met: %v", err)
	}
}

// SetupTestDB creates a test database with mock
func SetupTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	t.Helper()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to create GORM DB: %v", err)
	}

	return gormDB, mock
}

// SetupTestRedis creates a test Redis client with miniredis
func SetupTestRedis(t *testing.T) (*redis.Client, *miniredis.Miniredis) {
	t.Helper()

	mr := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	return client, mr
}

