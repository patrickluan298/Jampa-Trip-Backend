package repository

import (
	"github.com/jampa_trip/internal/app/model"
	"github.com/jampa_trip/internal/app/query"
	"gorm.io/gorm"
)

// FornecedorRepository - objeto de contexto
type FornecedorRepository struct {
	DB *gorm.DB
}

// FornecedorRepositoryNew - construtor do objeto
func FornecedorRepositoryNew(DB *gorm.DB) *FornecedorRepository {
	return &FornecedorRepository{
		DB: DB,
	}
}

// Get - busca um usuário pelo email
func (receiver *FornecedorRepository) Get(email string) (*model.Fornecedor, error) {
	row := &model.Fornecedor{}

	err := receiver.DB.Raw(query.ObterPorEmail, email).Row().Scan(
		&row.ID,
		&row.Nome,
		&row.Email,
		&row.Senha,
		&row.CNPJ,
		&row.Telefone,
		&row.Endereco,
		&row.MomentoCadastro,
	)

	if err != nil {
		return nil, err
	}

	return row, nil
}

// Cadastrar - cria um novo usuário
func (receiver *FornecedorRepository) Cadastrar(user *model.Fornecedor) error {
	err := receiver.DB.Raw(query.Cadastrar,
		user.Nome,
		user.Email,
		user.Senha,
		user.CNPJ,
		user.Telefone,
		user.Endereco,
		user.MomentoCadastro,
	).Row().Scan(&user.ID)

	return err
}

// EmailExiste - verifica se o email já está cadastrado
func (r *FornecedorRepository) EmailExiste(email string) (bool, error) {
	var count int64
	err := r.DB.Model(&model.Fornecedor{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

// UpdateUser - atualiza um usuário existente
func (r *FornecedorRepository) UpdateUser(user *model.Fornecedor) error {
	return r.DB.Save(user).Error
}

// DeleteUser - remove um usuário
func (r *FornecedorRepository) DeleteUser(id uint) error {
	return r.DB.Delete(&model.Fornecedor{}, id).Error
}
