package service

import (
	"net/http"
	"time"

	"github.com/jampa_trip/internal/app/contract"
	"github.com/jampa_trip/internal/app/model"
	"github.com/jampa_trip/internal/app/repository"
	"github.com/jampa_trip/internal/pkg/util"
	"gorm.io/gorm"
)

// ClienteService - objeto de contexto
type ClienteService struct {
	ClienteRepository *repository.ClienteRepository
}

// ClienteServiceNew - construtor do objeto
func ClienteServiceNew(DB *gorm.DB) *ClienteService {
	return &ClienteService{
		ClienteRepository: repository.ClienteRepositoryNew(DB),
	}
}

// Login - realiza a autenticação de um cliente
func (receiver *ClienteService) Login(request *contract.LoginClienteRequest) (*contract.LoginClienteResponse, error) {

	cliente, err := receiver.ClienteRepository.GetByEmail(request.Email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, util.WrapError("Email e/ou senha incorretos", nil, http.StatusUnauthorized)
		}
		return nil, util.WrapError("Erro ao buscar cliente", err, http.StatusInternalServerError)
	}

	if !util.VerificaSenha(request.Senha, cliente.Senha) {
		return nil, util.WrapError("Email e/ou senha incorretos", nil, http.StatusUnauthorized)
	}

	response := &contract.LoginClienteResponse{
		Mensagem: "Login realizado com sucesso",
		Dados: contract.ClienteLogin{
			ID:    cliente.ID,
			Nome:  cliente.Nome,
			Email: cliente.Email,
		},
	}

	return response, nil
}

// Create - realiza o cadastro de um novo cliente
func (receiver *ClienteService) Create(request *contract.CadastrarClienteRequest) (*contract.CadastrarClienteResponse, error) {

	emailExiste, err := receiver.ClienteRepository.EmailExiste(request.Email)
	if err != nil {
		return nil, util.WrapError("Erro ao verificar email", err, http.StatusInternalServerError)
	}

	if emailExiste {
		return nil, util.WrapError("O email informado já está cadastrado", nil, http.StatusConflict)
	}

	senhaHash, err := util.CriptografarSenha(request.Senha)
	if err != nil {
		return nil, util.WrapError("Erro ao criptografar senha", err, http.StatusInternalServerError)
	}

	dataNascimento, err := time.Parse("2006-01-02", request.DataNascimento)
	if err != nil {
		return nil, util.WrapError("Formato de data inválido. Use YYYY-MM-DD", err, http.StatusUnprocessableEntity)
	}

	cliente := &model.Cliente{
		Nome:               request.Nome,
		Email:              request.Email,
		Senha:              senhaHash,
		CPF:                request.CPF,
		Telefone:           request.Telefone,
		DataNascimento:     dataNascimento,
		MomentoCadastro:    time.Now(),
		MomentoAtualizacao: time.Now(),
	}

	if err := receiver.ClienteRepository.Create(cliente); err != nil {
		return nil, util.WrapError("Erro ao cadastrar cliente", err, http.StatusInternalServerError)
	}

	response := &contract.CadastrarClienteResponse{
		Mensagem: "Cliente cadastrado com sucesso",
		Dados: contract.Cliente{
			ID:                 cliente.ID,
			Nome:               cliente.Nome,
			Email:              cliente.Email,
			CPF:                cliente.CPF,
			Telefone:           cliente.Telefone,
			DataNascimento:     cliente.DataNascimento.Format("2006-01-02"),
			MomentoCadastro:    cliente.MomentoCadastro.Format("2006-01-02 15:04:05"),
			MomentoAtualizacao: cliente.MomentoAtualizacao.Format("2006-01-02 15:04:05"),
		},
	}

	return response, nil
}

// Update - realiza a atualização de um cliente existente
func (receiver *ClienteService) Update(request *contract.AtualizarClienteRequest) (*contract.AtualizarClienteResponse, error) {

	clienteExistente, err := receiver.ClienteRepository.GetByID(request.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, util.WrapError("Cliente não encontrado", nil, http.StatusNotFound)
		}
		return nil, util.WrapError("Erro ao buscar cliente", err, http.StatusInternalServerError)
	}

	emailExiste, err := receiver.ClienteRepository.EmailExisteParaOutroCliente(request.Email, request.ID)
	if err != nil {
		return nil, util.WrapError("Erro ao verificar email", err, http.StatusInternalServerError)
	}

	if emailExiste {
		return nil, util.WrapError("O email informado já está cadastrado para outro cliente", nil, http.StatusConflict)
	}

	senhaHash, err := util.CriptografarSenha(request.Senha)
	if err != nil {
		return nil, util.WrapError("Erro ao criptografar senha", err, http.StatusInternalServerError)
	}

	dataNascimento, err := time.Parse("2006-01-02", request.DataNascimento)
	if err != nil {
		return nil, util.WrapError("Formato de data inválido. Use YYYY-MM-DD", err, http.StatusUnprocessableEntity)
	}

	cliente := &model.Cliente{
		ID:                 request.ID,
		Nome:               request.Nome,
		Email:              request.Email,
		Senha:              senhaHash,
		CPF:                request.CPF,
		Telefone:           request.Telefone,
		DataNascimento:     dataNascimento,
		MomentoCadastro:    clienteExistente.MomentoCadastro,
		MomentoAtualizacao: time.Now(),
	}

	if err := receiver.ClienteRepository.Update(cliente); err != nil {
		return nil, util.WrapError("Erro ao atualizar cliente", err, http.StatusInternalServerError)
	}

	response := &contract.AtualizarClienteResponse{
		Mensagem: "Cliente atualizado com sucesso",
		Dados: contract.Cliente{
			ID:                 cliente.ID,
			Nome:               cliente.Nome,
			Email:              cliente.Email,
			CPF:                cliente.CPF,
			Telefone:           cliente.Telefone,
			DataNascimento:     cliente.DataNascimento.Format("2006-01-02"),
			MomentoCadastro:    cliente.MomentoCadastro.Format("2006-01-02 15:04:05"),
			MomentoAtualizacao: cliente.MomentoAtualizacao.Format("2006-01-02 15:04:05"),
		},
	}

	return response, nil
}

// List - realiza a listagem de todos os clientes
func (receiver *ClienteService) List() (*contract.ListarClienteResponse, error) {

	clientes, err := receiver.ClienteRepository.List()
	if err != nil {
		return nil, util.WrapError("Erro ao buscar clientes", err, http.StatusInternalServerError)
	}

	var clientesResponse []contract.Cliente
	for _, cliente := range clientes {
		clientesResponse = append(clientesResponse, contract.Cliente{
			ID:                 cliente.ID,
			Nome:               cliente.Nome,
			Email:              cliente.Email,
			CPF:                cliente.CPF,
			Telefone:           cliente.Telefone,
			DataNascimento:     cliente.DataNascimento.Format("2006-01-02"),
			MomentoCadastro:    cliente.MomentoCadastro.Format("2006-01-02 15:04:05"),
			MomentoAtualizacao: cliente.MomentoAtualizacao.Format("2006-01-02 15:04:05"),
		})
	}

	response := &contract.ListarClienteResponse{
		Mensagem: "Clientes listados com sucesso",
		Dados:    clientesResponse,
	}

	return response, nil
}

// Get - realiza a busca de um cliente por ID
func (receiver *ClienteService) Get(ID int) (*contract.ObterClienteResponse, error) {

	cliente, err := receiver.ClienteRepository.GetByID(ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, util.WrapError("Cliente não encontrado", nil, http.StatusNotFound)
		}
		return nil, util.WrapError("Erro ao buscar cliente", err, http.StatusInternalServerError)
	}

	response := &contract.ObterClienteResponse{
		Mensagem: "Cliente obtido com sucesso",
		Dados: contract.Cliente{
			ID:                 cliente.ID,
			Nome:               cliente.Nome,
			Email:              cliente.Email,
			CPF:                cliente.CPF,
			Telefone:           cliente.Telefone,
			DataNascimento:     cliente.DataNascimento.Format("2006-01-02"),
			MomentoCadastro:    cliente.MomentoCadastro.Format("2006-01-02 15:04:05"),
			MomentoAtualizacao: cliente.MomentoAtualizacao.Format("2006-01-02 15:04:05"),
		},
	}

	return response, nil
}
