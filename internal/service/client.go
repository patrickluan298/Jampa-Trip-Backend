package service

import (
	"net/http"
	"time"

	"github.com/jampa_trip/internal/contract"
	"github.com/jampa_trip/internal/model"
	"github.com/jampa_trip/internal/repository"
	"github.com/jampa_trip/pkg/util"
	"gorm.io/gorm"
)

// ClientService - objeto de contexto
type ClientService struct {
	ClientRepository *repository.ClientRepository
}

// ClientServiceNew - construtor do objeto
func ClientServiceNew(DB *gorm.DB) *ClientService {
	return &ClientService{
		ClientRepository: repository.ClientRepositoryNew(DB),
	}
}

// Create - realiza o cadastro de um novo cliente
func (receiver *ClientService) Create(request *contract.CreateClientRequest) (*contract.CreateClientResponse, error) {

	if request.Password != request.ConfirmPassword {
		return nil, util.WrapError("As senhas não coincidem", nil, http.StatusUnprocessableEntity)
	}

	emailExists, err := receiver.ClientRepository.EmailExiste(request.Email)
	if err != nil {
		return nil, util.WrapError("Erro ao verificar email", err, http.StatusInternalServerError)
	}

	if emailExists {
		return nil, util.WrapError("O email informado já está cadastrado", nil, http.StatusConflict)
	}

	passwordHash, err := util.CriptografarSenha(request.Password)
	if err != nil {
		return nil, util.WrapError("Erro ao criptografar senha", err, http.StatusInternalServerError)
	}

	birthDate, err := time.Parse("2006-01-02", request.BirthDate)
	if err != nil {
		return nil, util.WrapError("Formato de data inválido. Use YYYY-MM-DD", err, http.StatusUnprocessableEntity)
	}

	client := &model.Client{
		Name:      request.Name,
		Email:     request.Email,
		Password:  passwordHash,
		CPF:       request.CPF,
		Phone:     request.Phone,
		BirthDate: birthDate,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := receiver.ClientRepository.Create(client); err != nil {
		return nil, util.WrapError("Erro ao cadastrar cliente", err, http.StatusInternalServerError)
	}

	response := &contract.CreateClientResponse{
		Message: "Cliente cadastrado com sucesso",
		Data: contract.Client{
			ID:        client.ID,
			Name:      client.Name,
			Email:     client.Email,
			CPF:       client.CPF,
			Phone:     client.Phone,
			BirthDate: client.BirthDate.Format("2006-01-02"),
			CreatedAt: client.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: client.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	}

	return response, nil
}

// Update - realiza a atualização de um cliente existente
func (receiver *ClientService) Update(request *contract.UpdateClientRequest) (*contract.UpdateClientResponse, error) {

	existingClient, err := receiver.ClientRepository.GetByID(request.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, util.WrapError("Cliente não encontrado", nil, http.StatusNotFound)
		}
		return nil, util.WrapError("Erro ao buscar cliente", err, http.StatusInternalServerError)
	}

	emailExists, err := receiver.ClientRepository.EmailExisteParaOutroCliente(request.Email, request.ID)
	if err != nil {
		return nil, util.WrapError("Erro ao verificar email", err, http.StatusInternalServerError)
	}

	if emailExists {
		return nil, util.WrapError("O email informado já está cadastrado para outro cliente", nil, http.StatusConflict)
	}

	passwordHash, err := util.CriptografarSenha(request.Password)
	if err != nil {
		return nil, util.WrapError("Erro ao criptografar senha", err, http.StatusInternalServerError)
	}

	birthDate, err := time.Parse("2006-01-02", request.BirthDate)
	if err != nil {
		return nil, util.WrapError("Formato de data inválido. Use YYYY-MM-DD", err, http.StatusUnprocessableEntity)
	}

	client := &model.Client{
		ID:        request.ID,
		Name:      request.Name,
		Email:     request.Email,
		Password:  passwordHash,
		CPF:       request.CPF,
		Phone:     request.Phone,
		BirthDate: birthDate,
		CreatedAt: existingClient.CreatedAt,
		UpdatedAt: time.Now(),
	}

	if err := receiver.ClientRepository.Update(client); err != nil {
		return nil, util.WrapError("Erro ao atualizar cliente", err, http.StatusInternalServerError)
	}

	response := &contract.UpdateClientResponse{
		Message: "Cliente atualizado com sucesso",
		Data: contract.Client{
			ID:        client.ID,
			Name:      client.Name,
			Email:     client.Email,
			CPF:       client.CPF,
			Phone:     client.Phone,
			BirthDate: client.BirthDate.Format("2006-01-02"),
			CreatedAt: client.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: client.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	}

	return response, nil
}

// List - realiza a listagem de todos os clientes
func (receiver *ClientService) List(filtros *model.Client) (*contract.ListClientResponse, error) {

	clients, err := receiver.ClientRepository.List(filtros)
	if err != nil {
		return nil, util.WrapError("Erro ao buscar clientes", err, http.StatusInternalServerError)
	}

	var clientsResponse []contract.Client
	for _, client := range clients {
		clientsResponse = append(clientsResponse, contract.Client{
			ID:        client.ID,
			Name:      client.Name,
			Email:     client.Email,
			CPF:       client.CPF,
			Phone:     client.Phone,
			BirthDate: client.BirthDate.Format("2006-01-02"),
			CreatedAt: client.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: client.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	response := &contract.ListClientResponse{
		Message: "Clientes listados com sucesso",
		Data:    clientsResponse,
	}

	return response, nil
}

// Get - realiza a busca de um cliente por ID
func (receiver *ClientService) Get(ID int) (*contract.GetClientResponse, error) {

	client, err := receiver.ClientRepository.GetByID(ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, util.WrapError("Cliente não encontrado", nil, http.StatusNotFound)
		}
		return nil, util.WrapError("Erro ao buscar cliente", err, http.StatusInternalServerError)
	}

	response := &contract.GetClientResponse{
		Message: "Cliente obtido com sucesso",
		Data: contract.Client{
			ID:        client.ID,
			Name:      client.Name,
			Email:     client.Email,
			CPF:       client.CPF,
			Phone:     client.Phone,
			BirthDate: client.BirthDate.Format("2006-01-02"),
			CreatedAt: client.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: client.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	}

	return response, nil
}
