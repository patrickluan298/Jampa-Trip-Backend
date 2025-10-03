package repository

import (
	"github.com/jampa_trip/internal/model"
	"github.com/jampa_trip/internal/query"
	"gorm.io/gorm"
)

// ClienteRepository - objeto de contexto
type ClienteRepository struct {
	DB *gorm.DB
}

// ClienteRepositoryNew - construtor do objeto
func ClienteRepositoryNew(DB *gorm.DB) *ClienteRepository {
	return &ClienteRepository{
		DB: DB,
	}
}

// GetByID - busca um cliente pelo ID
func (receiver *ClienteRepository) GetByID(id int) (*model.Cliente, error) {
	row := &model.Cliente{}

	err := receiver.DB.Raw(query.ObterClientePorID, id).Row().Scan(
		&row.ID,
		&row.Nome,
		&row.Email,
		&row.Senha,
		&row.CPF,
		&row.Telefone,
		&row.DataNascimento,
		&row.MomentoCadastro,
		&row.MomentoAtualizacao,
	)

	if err != nil {
		return nil, err
	}

	return row, nil
}

// GetByEmail - busca um cliente pelo email
func (receiver *ClienteRepository) GetByEmail(email string) (*model.Cliente, error) {
	row := &model.Cliente{}

	err := receiver.DB.Raw(query.ObterClientePorEmail, email).Row().Scan(
		&row.ID,
		&row.Nome,
		&row.Email,
		&row.Senha,
		&row.CPF,
		&row.Telefone,
		&row.DataNascimento,
		&row.MomentoCadastro,
		&row.MomentoAtualizacao,
	)

	if err != nil {
		return nil, err
	}

	return row, nil
}

// Create - cria um novo cliente
func (receiver *ClienteRepository) Create(cliente *model.Cliente) error {
	err := receiver.DB.Raw(query.CadastrarCliente,
		cliente.Nome,
		cliente.Email,
		cliente.Senha,
		cliente.CPF,
		cliente.Telefone,
		cliente.DataNascimento,
		cliente.MomentoCadastro,
		cliente.MomentoAtualizacao,
	).Row().Scan(&cliente.ID)

	return err
}

// Update - atualiza um cliente existente
func (receiver *ClienteRepository) Update(cliente *model.Cliente) error {
	err := receiver.DB.Raw(query.AtualizarCliente,
		cliente.Nome,
		cliente.Email,
		cliente.Senha,
		cliente.CPF,
		cliente.Telefone,
		cliente.DataNascimento,
		cliente.MomentoAtualizacao,
		cliente.ID,
	).Row().Scan()

	return err
}

// List - busca todos os clientes
func (receiver *ClienteRepository) List(filtros *model.Cliente) ([]*model.Cliente, error) {
	rows, err := receiver.DB.Raw(query.ListarTodosClientes, filtros.Nome, filtros.Nome, filtros.Email, filtros.Email,
		filtros.CPF, filtros.CPF, filtros.Telefone, filtros.Telefone, filtros.DataNascimento, filtros.DataNascimento).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clientes []*model.Cliente
	for rows.Next() {
		cliente := &model.Cliente{}
		err := rows.Scan(
			&cliente.ID,
			&cliente.Nome,
			&cliente.Email,
			&cliente.CPF,
			&cliente.Telefone,
			&cliente.DataNascimento,
			&cliente.MomentoCadastro,
			&cliente.MomentoAtualizacao,
		)
		if err != nil {
			return nil, err
		}
		clientes = append(clientes, cliente)
	}

	return clientes, nil
}

// EmailExiste - verifica se o email já está cadastrado
func (r *ClienteRepository) EmailExiste(email string) (bool, error) {
	var count int64
	err := r.DB.Model(&model.Cliente{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

// EmailExisteParaOutroCliente - verifica se o email já está cadastrado para outro cliente
func (r *ClienteRepository) EmailExisteParaOutroCliente(email string, id int) (bool, error) {
	var count int64
	err := r.DB.Model(&model.Cliente{}).Where("email = ? AND id != ?", email, id).Count(&count).Error
	return count > 0, err
}
