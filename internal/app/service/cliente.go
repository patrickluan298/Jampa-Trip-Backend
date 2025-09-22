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

	token, err := util.GenerateToken()
	if err != nil {
		return nil, util.WrapError("Erro ao gerar token", err, http.StatusInternalServerError)
	}

	response := &contract.LoginClienteResponse{
		Mensagem: "Login realizado com sucesso",
		Token:    token,
		Dados: contract.ClienteLogin{
			ID:    cliente.ID,
			Nome:  cliente.Nome,
			Email: cliente.Email,
		},
	}

	return response, nil
}

// Cadastrar - realiza o cadastro de um novo cliente
func (receiver *ClienteService) Cadastrar(request *contract.CadastrarClienteRequest) (*contract.CadastrarClienteResponse, error) {

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

	if err := receiver.ClienteRepository.Cadastrar(cliente); err != nil {
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

// Atualizar - realiza a atualização de um cliente existente
func (receiver *ClienteService) Atualizar(request *contract.AtualizarClienteRequest) (*contract.AtualizarClienteResponse, error) {

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

	if err := receiver.ClienteRepository.Atualizar(cliente); err != nil {
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

// Listar - realiza a listagem de todos os clientes
func (receiver *ClienteService) Listar() (*contract.ListarClienteResponse, error) {

	clientes, err := receiver.ClienteRepository.ListarTodos()
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
