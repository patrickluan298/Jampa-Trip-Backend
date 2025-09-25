package repository

import (
	"github.com/jampa_trip/internal/app/model"
	"github.com/jampa_trip/internal/app/query"
	"gorm.io/gorm"
)

// EmpresaRepository - objeto de contexto
type EmpresaRepository struct {
	DB *gorm.DB
}

// EmpresaRepositoryNew - construtor do objeto
func EmpresaRepositoryNew(DB *gorm.DB) *EmpresaRepository {
	return &EmpresaRepository{
		DB: DB,
	}
}

// GetByID - busca uma empresa pelo ID
func (receiver *EmpresaRepository) GetByID(id int) (*model.Empresa, error) {
	row := &model.Empresa{}

	err := receiver.DB.Raw(query.ObterEmpresaPorID, id).Row().Scan(
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

// GetByEmail - busca uma empresa pelo email
func (receiver *EmpresaRepository) GetByEmail(email string) (*model.Empresa, error) {
	row := &model.Empresa{}

	err := receiver.DB.Raw(query.ObterEmpresaPorEmail, email).Row().Scan(
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

// Create - cria uma nova empresa
func (receiver *EmpresaRepository) Create(empresa *model.Empresa) error {
	err := receiver.DB.Raw(query.CadastrarEmpresa,
		empresa.Nome,
		empresa.Email,
		empresa.Senha,
		empresa.CNPJ,
		empresa.Telefone,
		empresa.Endereco,
		empresa.MomentoCadastro,
		empresa.MomentoAtualizacao,
	).Row().Scan(&empresa.ID)

	return err
}

// Update - atualiza uma empresa existente
func (receiver *EmpresaRepository) Update(empresa *model.Empresa) error {
	err := receiver.DB.Raw(query.AtualizarEmpresa,
		empresa.Nome,
		empresa.Email,
		empresa.Senha,
		empresa.CNPJ,
		empresa.Telefone,
		empresa.Endereco,
		empresa.MomentoAtualizacao,
		empresa.ID,
	).Row().Scan()

	return err
}

// List - busca todas as empresas
func (receiver *EmpresaRepository) List() ([]*model.Empresa, error) {
	rows, err := receiver.DB.Raw(query.ListarTodasEmpresas).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var empresas []*model.Empresa
	for rows.Next() {
		empresa := &model.Empresa{}
		err := rows.Scan(
			&empresa.ID,
			&empresa.Nome,
			&empresa.Email,
			&empresa.CNPJ,
			&empresa.Telefone,
			&empresa.Endereco,
			&empresa.MomentoCadastro,
			&empresa.MomentoAtualizacao,
		)
		if err != nil {
			return nil, err
		}
		empresas = append(empresas, empresa)
	}

	return empresas, nil
}

// EmailExiste - verifica se o email já está cadastrado
func (r *EmpresaRepository) EmailExiste(email string) (bool, error) {
	var count int64
	err := r.DB.Model(&model.Empresa{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

// EmailExisteParaOutraEmpresa - verifica se o email já está cadastrado para outra empresa
func (r *EmpresaRepository) EmailExisteParaOutraEmpresa(email string, id int) (bool, error) {
	var count int64
	err := r.DB.Model(&model.Empresa{}).Where("email = ? AND id != ?", email, id).Count(&count).Error
	return count > 0, err
}
