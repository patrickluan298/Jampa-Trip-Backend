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

// FornecedorService - objeto de contexto
type FornecedorService struct {
	FornecedorRepository *repository.FornecedorRepository
}

// FornecedorServiceNew - construtor do objeto
func FornecedorServiceNew(DB *gorm.DB) *FornecedorService {
	return &FornecedorService{
		FornecedorRepository: repository.FornecedorRepositoryNew(DB),
	}
}

// Login - realiza a autenticação de um usuário
func (receiver *FornecedorService) Login(request *contract.LoginFornecedorRequest) (*contract.LoginFornecedorResponse, error) {

	fornecedor, err := receiver.FornecedorRepository.GetByEmail(request.Email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, util.WrapError("Email e/ou senha incorretos", nil, http.StatusUnauthorized)
		}
		return nil, util.WrapError("Erro ao buscar usuário", err, http.StatusInternalServerError)
	}

	if !util.VerificaSenha(request.Senha, fornecedor.Senha) {
		return nil, util.WrapError("Email e/ou senha incorretos", nil, http.StatusUnauthorized)
	}

	token, err := util.GenerateToken()
	if err != nil {
		return nil, util.WrapError("Erro ao gerar token", err, http.StatusInternalServerError)
	}

	response := &contract.LoginFornecedorResponse{
		Mensagem: "Login realizado com sucesso",
		Token:    token,
		Dados: contract.Fornecedor{
			ID:              fornecedor.ID,
			Nome:            fornecedor.Nome,
			Email:           fornecedor.Email,
			CNPJ:            fornecedor.CNPJ,
			Telefone:        fornecedor.Telefone,
			Endereco:        fornecedor.Endereco,
			MomentoCadastro: fornecedor.MomentoCadastro.Format(time.RFC3339),
		},
	}

	return response, nil
}

// Cadastrar - realiza o cadastro de um novo fornecedor
func (receiver *FornecedorService) Cadastrar(request *contract.CadastrarFornecedorRequest) (*contract.CadastrarFornecedorResponse, error) {

	emailExiste, err := receiver.FornecedorRepository.EmailExiste(request.Email)
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

	fornecedor := &model.Fornecedor{
		Nome:               request.Nome,
		Email:              request.Email,
		Senha:              senhaHash,
		CNPJ:               request.CNPJ,
		Telefone:           request.Telefone,
		Endereco:           request.Endereco,
		MomentoCadastro:    time.Now(),
		MomentoAtualizacao: time.Now(),
	}

	if err := receiver.FornecedorRepository.Cadastrar(fornecedor); err != nil {
		return nil, util.WrapError("Erro ao cadastrar fornecedor", err, http.StatusInternalServerError)
	}

	response := &contract.CadastrarFornecedorResponse{
		Mensagem: "Fornecedor cadastrado com sucesso",
		Dados: contract.Fornecedor{
			ID:                 fornecedor.ID,
			Nome:               fornecedor.Nome,
			Email:              fornecedor.Email,
			CNPJ:               fornecedor.CNPJ,
			Telefone:           fornecedor.Telefone,
			Endereco:           fornecedor.Endereco,
			MomentoCadastro:    fornecedor.MomentoCadastro.Format("2006-01-02 15:04:05"),
			MomentoAtualizacao: fornecedor.MomentoAtualizacao.Format("2006-01-02 15:04:05"),
		},
	}

	return response, nil
}

// Atualizar - realiza a atualização de um fornecedor existente
func (receiver *FornecedorService) Atualizar(request *contract.AtualizarFornecedorRequest) (*contract.AtualizarFornecedorResponse, error) {

	fornecedorExistente, err := receiver.FornecedorRepository.GetByID(request.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, util.WrapError("Fornecedor não encontrado", nil, http.StatusNotFound)
		}
		return nil, util.WrapError("Erro ao buscar fornecedor", err, http.StatusInternalServerError)
	}

	emailExiste, err := receiver.FornecedorRepository.EmailExisteParaOutroFornecedor(request.Email, request.ID)
	if err != nil {
		return nil, util.WrapError("Erro ao verificar email", err, http.StatusInternalServerError)
	}

	if emailExiste {
		return nil, util.WrapError("O email informado já está cadastrado para outro fornecedor", nil, http.StatusConflict)
	}

	senhaHash, err := util.CriptografarSenha(request.Senha)
	if err != nil {
		return nil, util.WrapError("Erro ao criptografar senha", err, http.StatusInternalServerError)
	}

	fornecedor := &model.Fornecedor{
		ID:                 request.ID,
		Nome:               request.Nome,
		Email:              request.Email,
		Senha:              senhaHash,
		CNPJ:               request.CNPJ,
		Telefone:           request.Telefone,
		Endereco:           request.Endereco,
		MomentoCadastro:    fornecedorExistente.MomentoCadastro,
		MomentoAtualizacao: time.Now(),
	}

	if err := receiver.FornecedorRepository.Atualizar(fornecedor); err != nil {
		return nil, util.WrapError("Erro ao atualizar fornecedor", err, http.StatusInternalServerError)
	}

	response := &contract.AtualizarFornecedorResponse{
		Mensagem: "Fornecedor atualizado com sucesso",
		Dados: contract.Fornecedor{
			ID:                 fornecedor.ID,
			Nome:               fornecedor.Nome,
			Email:              fornecedor.Email,
			CNPJ:               fornecedor.CNPJ,
			Telefone:           fornecedor.Telefone,
			Endereco:           fornecedor.Endereco,
			MomentoCadastro:    fornecedor.MomentoCadastro.Format("2006-01-02 15:04:05"),
			MomentoAtualizacao: fornecedor.MomentoAtualizacao.Format("2006-01-02 15:04:05"),
		},
	}

	return response, nil
}
