package repository

import (
	"github.com/jampa_trip/internal/app/model"
	"gorm.io/gorm"
)

// AuthRepository objeto de contexto
type FornecedorRepository struct {
	DB *gorm.DB
}

// FornecedorRepositoryNew construtor do objeto
func FornecedorRepositoryNew(DB *gorm.DB) *FornecedorRepository {
	return &FornecedorRepository{
		DB: DB,
	}
}

// GetFornecedor busca um usuário pelo email
func (r *FornecedorRepository) GetFornecedor(email string) (*model.Fornecedor, error) {
	var user model.Fornecedor
	err := r.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// CreateUser cria um novo usuário
func (r *FornecedorRepository) CreateUser(user *model.Fornecedor) error {
	return r.DB.Create(user).Error
}

// UpdateUser atualiza um usuário existente
func (r *FornecedorRepository) UpdateUser(user *model.Fornecedor) error {
	return r.DB.Save(user).Error
}

// DeleteUser remove um usuário
func (r *FornecedorRepository) DeleteUser(id uint) error {
	return r.DB.Delete(&model.Fornecedor{}, id).Error
}
