package repository

import (
	"github.com/jampa_trip/internal/model"
	"github.com/jampa_trip/internal/query"
	"gorm.io/gorm"
)

// ClientRepository - objeto de contexto
type ClientRepository struct {
	DB *gorm.DB
}

// ClientRepositoryNew - construtor do objeto
func ClientRepositoryNew(DB *gorm.DB) *ClientRepository {
	return &ClientRepository{
		DB: DB,
	}
}

// GetByID - busca um cliente pelo ID
func (receiver *ClientRepository) GetByID(id int) (*model.Client, error) {
	row := &model.Client{}

	err := receiver.DB.Raw(query.GetClientByID, id).Row().Scan(
		&row.ID,
		&row.Name,
		&row.Email,
		&row.Password,
		&row.CPF,
		&row.Phone,
		&row.BirthDate,
		&row.CreatedAt,
		&row.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return row, nil
}

// GetByEmail - busca um cliente pelo email
func (receiver *ClientRepository) GetByEmail(email string) (*model.Client, error) {
	row := &model.Client{}

	err := receiver.DB.Raw(query.GetClientByEmail, email).Row().Scan(
		&row.ID,
		&row.Name,
		&row.Email,
		&row.Password,
		&row.CPF,
		&row.Phone,
		&row.BirthDate,
		&row.CreatedAt,
		&row.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return row, nil
}

// Create - cria um novo cliente
func (receiver *ClientRepository) Create(client *model.Client) error {
	err := receiver.DB.Raw(query.CreateClient,
		client.Name,
		client.Email,
		client.Password,
		client.CPF,
		client.Phone,
		client.BirthDate,
		client.CreatedAt,
		client.UpdatedAt,
	).Row().Scan(&client.ID)

	return err
}

// Update - atualiza os campos enviados no map
func (receiver *ClientRepository) Update(id int, updates map[string]interface{}) error {
	result := receiver.DB.Model(&model.Client{}).Where("id = ?", id).Updates(updates)
	return result.Error
}

// List - busca todos os clientes
func (receiver *ClientRepository) List(filtros *model.Client) ([]*model.Client, error) {
	rows, err := receiver.DB.Raw(query.ListAllClients, filtros.Name, filtros.Name, filtros.Email, filtros.Email,
		filtros.CPF, filtros.CPF, filtros.Phone, filtros.Phone, filtros.BirthDate, filtros.BirthDate).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clients []*model.Client
	for rows.Next() {
		client := &model.Client{}
		err := rows.Scan(
			&client.ID,
			&client.Name,
			&client.Email,
			&client.CPF,
			&client.Phone,
			&client.BirthDate,
			&client.CreatedAt,
			&client.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		clients = append(clients, client)
	}

	return clients, nil
}

// EmailExiste - verifica se o email já está cadastrado
func (r *ClientRepository) EmailExiste(email string) (bool, error) {
	var count int64
	err := r.DB.Model(&model.Client{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

// EmailExisteParaOutroCliente - verifica se o email já está cadastrado para outro cliente
func (r *ClientRepository) EmailExisteParaOutroCliente(email string, id int) (bool, error) {
	var count int64
	err := r.DB.Model(&model.Client{}).Where("email = ? AND id != ?", email, id).Count(&count).Error
	return count > 0, err
}
