package repository

import "gorm.io/gorm"

// AuthRepository objeto de contexto
type AuthRepository struct {
	DB *gorm.DB
}

// AuthRepositoryNew construtor do objeto
func AuthRepositoryNew(DB *gorm.DB) *AuthRepository {
	return &AuthRepository{
		DB: DB,
	}
}
