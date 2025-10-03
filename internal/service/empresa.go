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

// EmpresaService - objeto de contexto
type EmpresaService struct {
	EmpresaRepository *repository.EmpresaRepository
}

// EmpresaServiceNew - construtor do objeto
func EmpresaServiceNew(DB *gorm.DB) *EmpresaService {
	return &EmpresaService{
		EmpresaRepository: repository.EmpresaRepositoryNew(DB),
	}
}

// Login - realiza a autenticação de uma empresa
func (receiver *EmpresaService) Login(request *contract.LoginEmpresaRequest) (*contract.LoginEmpresaResponse, error) {

	empresa, err := receiver.EmpresaRepository.GetByEmail(request.Email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, util.WrapError("Email e/ou senha incorretos", nil, http.StatusUnauthorized)
		}
		return nil, util.WrapError("Erro ao buscar empresa", err, http.StatusInternalServerError)
	}

	if !util.VerificaSenha(request.Senha, empresa.Senha) {
		return nil, util.WrapError("Email e/ou senha incorretos", nil, http.StatusUnauthorized)
	}

	response := &contract.LoginEmpresaResponse{
		Mensagem: "Login realizado com sucesso",
		Dados: contract.EmpresaLogin{
			ID:    empresa.ID,
			Nome:  empresa.Nome,
			Email: empresa.Email,
		},
	}

	return response, nil
}

// Create - realiza o cadastro de uma nova empresa
func (receiver *EmpresaService) Create(request *contract.CadastrarEmpresaRequest) (*contract.CadastrarEmpresaResponse, error) {

	if request.Senha != request.ConfirmarSenha {
		return nil, util.WrapError("As senhas não coincidem", nil, http.StatusUnprocessableEntity)
	}

	emailExiste, err := receiver.EmpresaRepository.EmailExiste(request.Email)
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

	empresa := &model.Empresa{
		Nome:               request.Nome,
		Email:              request.Email,
		Senha:              senhaHash,
		CNPJ:               request.CNPJ,
		Telefone:           request.Telefone,
		Endereco:           request.Endereco,
		MomentoCadastro:    time.Now(),
		MomentoAtualizacao: time.Now(),
	}

	if err := receiver.EmpresaRepository.Create(empresa); err != nil {
		return nil, util.WrapError("Erro ao cadastrar empresa", err, http.StatusInternalServerError)
	}

	response := &contract.CadastrarEmpresaResponse{
		Mensagem: "Empresa cadastrada com sucesso",
		Dados: contract.Empresa{
			ID:                 empresa.ID,
			Nome:               empresa.Nome,
			Email:              empresa.Email,
			CNPJ:               empresa.CNPJ,
			Telefone:           empresa.Telefone,
			Endereco:           empresa.Endereco,
			MomentoCadastro:    empresa.MomentoCadastro.Format("2006-01-02 15:04:05"),
			MomentoAtualizacao: empresa.MomentoAtualizacao.Format("2006-01-02 15:04:05"),
		},
	}

	return response, nil
}

// Update - realiza a atualização de uma empresa existente
func (receiver *EmpresaService) Update(request *contract.AtualizarEmpresaRequest) (*contract.AtualizarEmpresaResponse, error) {

	empresaExistente, err := receiver.EmpresaRepository.GetByID(request.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, util.WrapError("Empresa não encontrada", nil, http.StatusNotFound)
		}
		return nil, util.WrapError("Erro ao buscar empresa", err, http.StatusInternalServerError)
	}

	emailExiste, err := receiver.EmpresaRepository.EmailExisteParaOutraEmpresa(request.Email, request.ID)
	if err != nil {
		return nil, util.WrapError("Erro ao verificar email", err, http.StatusInternalServerError)
	}

	if emailExiste {
		return nil, util.WrapError("O email informado já está cadastrado para outra empresa", nil, http.StatusConflict)
	}

	senhaHash, err := util.CriptografarSenha(request.Senha)
	if err != nil {
		return nil, util.WrapError("Erro ao criptografar senha", err, http.StatusInternalServerError)
	}

	empresa := &model.Empresa{
		ID:                 request.ID,
		Nome:               request.Nome,
		Email:              request.Email,
		Senha:              senhaHash,
		CNPJ:               request.CNPJ,
		Telefone:           request.Telefone,
		Endereco:           request.Endereco,
		MomentoCadastro:    empresaExistente.MomentoCadastro,
		MomentoAtualizacao: time.Now(),
	}

	if err := receiver.EmpresaRepository.Update(empresa); err != nil {
		return nil, util.WrapError("Erro ao atualizar empresa", err, http.StatusInternalServerError)
	}

	response := &contract.AtualizarEmpresaResponse{
		Mensagem: "Empresa atualizada com sucesso",
		Dados: contract.Empresa{
			ID:                 empresa.ID,
			Nome:               empresa.Nome,
			Email:              empresa.Email,
			CNPJ:               empresa.CNPJ,
			Telefone:           empresa.Telefone,
			Endereco:           empresa.Endereco,
			MomentoCadastro:    empresa.MomentoCadastro.Format("2006-01-02 15:04:05"),
			MomentoAtualizacao: empresa.MomentoAtualizacao.Format("2006-01-02 15:04:05"),
		},
	}

	return response, nil
}

// List - realiza a listagem de todas as empresas
func (receiver *EmpresaService) List(filtros *model.Empresa) (*contract.ListarEmpresasResponse, error) {

	empresas, err := receiver.EmpresaRepository.List(filtros)
	if err != nil {
		return nil, util.WrapError("Erro ao buscar empresas", err, http.StatusInternalServerError)
	}

	var empresasResponse []contract.Empresa
	for _, empresa := range empresas {
		empresasResponse = append(empresasResponse, contract.Empresa{
			ID:                 empresa.ID,
			Nome:               empresa.Nome,
			Email:              empresa.Email,
			CNPJ:               empresa.CNPJ,
			Telefone:           empresa.Telefone,
			Endereco:           empresa.Endereco,
			MomentoCadastro:    empresa.MomentoCadastro.Format("2006-01-02 15:04:05"),
			MomentoAtualizacao: empresa.MomentoAtualizacao.Format("2006-01-02 15:04:05"),
		})
	}

	response := &contract.ListarEmpresasResponse{
		Mensagem: "Empresas listadas com sucesso",
		Dados:    empresasResponse,
	}

	return response, nil
}

// Get - realiza a busca de uma empresa por ID
func (receiver *EmpresaService) Get(ID int) (*contract.ObterEmpresaResponse, error) {

	empresa, err := receiver.EmpresaRepository.GetByID(ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, util.WrapError("Empresa não encontrada", nil, http.StatusNotFound)
		}
		return nil, util.WrapError("Erro ao buscar empresa", err, http.StatusInternalServerError)
	}

	response := &contract.ObterEmpresaResponse{
		Mensagem: "Empresa obtida com sucesso",
		Dados: contract.Empresa{
			ID:                 empresa.ID,
			Nome:               empresa.Nome,
			Email:              empresa.Email,
			CNPJ:               empresa.CNPJ,
			Telefone:           empresa.Telefone,
			Endereco:           empresa.Endereco,
			MomentoCadastro:    empresa.MomentoCadastro.Format("2006-01-02 15:04:05"),
			MomentoAtualizacao: empresa.MomentoAtualizacao.Format("2006-01-02 15:04:05"),
		},
	}

	return response, nil
}
