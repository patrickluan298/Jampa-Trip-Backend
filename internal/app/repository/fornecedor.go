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

// GetByID - busca um fornecedor pelo ID
func (receiver *FornecedorRepository) GetByID(id int) (*model.Fornecedor, error) {
	row := &model.Fornecedor{}

	err := receiver.DB.Raw(query.ObterPorID, id).Row().Scan(
		&row.ID,
		&row.Nome,
		&row.Email,
		&row.Senha,
		&row.CNPJ,
		&row.Telefone,
		&row.Endereco,
		&row.MomentoCadastro,
		&row.MomentoAtualizacao,
	)

	if err != nil {
		return nil, err
	}

	return row, nil
}

// GetByEmail - busca um fornecedor pelo email
func (receiver *FornecedorRepository) GetByEmail(email string) (*model.Fornecedor, error) {
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
		&row.MomentoAtualizacao,
	)

	if err != nil {
		return nil, err
	}

	return row, nil
}

// Cadastrar - cria um novo fornecedor
func (receiver *FornecedorRepository) Cadastrar(fornecedor *model.Fornecedor) error {
	err := receiver.DB.Raw(query.Cadastrar,
		fornecedor.Nome,
		fornecedor.Email,
		fornecedor.Senha,
		fornecedor.CNPJ,
		fornecedor.Telefone,
		fornecedor.Endereco,
		fornecedor.MomentoCadastro,
		fornecedor.MomentoAtualizacao,
	).Row().Scan(&fornecedor.ID)

	return err
}

// Atualizar - atualiza um fornecedor existente
func (receiver *FornecedorRepository) Atualizar(fornecedor *model.Fornecedor) error {
	err := receiver.DB.Raw(query.Atualizar,
		fornecedor.Nome,
		fornecedor.Email,
		fornecedor.Senha,
		fornecedor.CNPJ,
		fornecedor.Telefone,
		fornecedor.Endereco,
		fornecedor.MomentoAtualizacao,
		fornecedor.ID,
	).Row().Scan()

	return err
}

// EmailExiste - verifica se o email já está cadastrado
func (r *FornecedorRepository) EmailExiste(email string) (bool, error) {
	var count int64
	err := r.DB.Model(&model.Fornecedor{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

// EmailExisteParaOutroFornecedor - verifica se o email já está cadastrado para outro fornecedor
func (r *FornecedorRepository) EmailExisteParaOutroFornecedor(email string, id int) (bool, error) {
	var count int64
	err := r.DB.Model(&model.Fornecedor{}).Where("email = ? AND id != ?", email, id).Count(&count).Error
	return count > 0, err
}
