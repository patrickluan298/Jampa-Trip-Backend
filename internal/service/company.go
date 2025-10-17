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

// CompanyService - objeto de contexto
type CompanyService struct {
	CompanyRepository *repository.CompanyRepository
}

// CompanyServiceNew - construtor do objeto
func CompanyServiceNew(DB *gorm.DB) *CompanyService {
	return &CompanyService{
		CompanyRepository: repository.CompanyRepositoryNew(DB),
	}
}

// Create - realiza o cadastro de uma nova empresa
func (receiver *CompanyService) Create(request *contract.CreateCompanyRequest) (*contract.CreateCompanyResponse, error) {

	if request.Password != request.ConfirmPassword {
		return nil, util.WrapError("As senhas não coincidem", nil, http.StatusUnprocessableEntity)
	}

	emailExists, err := receiver.CompanyRepository.EmailExiste(request.Email)
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

	company := &model.Company{
		Name:      request.Name,
		Email:     request.Email,
		Password:  passwordHash,
		CNPJ:      request.CNPJ,
		Phone:     request.Phone,
		Address:   request.Address,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := receiver.CompanyRepository.Create(company); err != nil {
		return nil, util.WrapError("Erro ao cadastrar empresa", err, http.StatusInternalServerError)
	}

	response := &contract.CreateCompanyResponse{
		Message: "Empresa cadastrada com sucesso",
		Data: contract.Company{
			ID:        company.ID,
			Name:      company.Name,
			Email:     company.Email,
			CNPJ:      company.CNPJ,
			Phone:     company.Phone,
			Address:   company.Address,
			CreatedAt: company.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: company.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	}

	return response, nil
}

// Update - realiza a atualização de uma empresa existente
func (receiver *CompanyService) Update(request *contract.UpdateCompanyRequest) (*contract.UpdateCompanyResponse, error) {

	if _, err := receiver.CompanyRepository.GetByID(request.ID); err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, util.WrapError("Empresa não encontrada", nil, http.StatusNotFound)
		}
		return nil, util.WrapError("Erro ao buscar empresa", err, http.StatusInternalServerError)
	}

	updates := make(map[string]interface{})

	if request.Name != nil {
		updates["name"] = *request.Name
	}

	if request.Email != nil {
		emailExists, err := receiver.CompanyRepository.EmailExisteParaOutraEmpresa(*request.Email, request.ID)
		if err != nil {
			return nil, util.WrapError("Erro ao verificar email", err, http.StatusInternalServerError)
		}

		if emailExists {
			return nil, util.WrapError("O email informado já está cadastrado para outra empresa", nil, http.StatusConflict)
		}
		updates["email"] = *request.Email
	}

	if request.Password != nil && request.ConfirmPassword != nil {
		if *request.Password != *request.ConfirmPassword {
			return nil, util.WrapError("As senhas não coincidem", nil, http.StatusUnprocessableEntity)
		}

		passwordHash, err := util.CriptografarSenha(*request.Password)
		if err != nil {
			return nil, util.WrapError("Erro ao criptografar senha", err, http.StatusInternalServerError)
		}
		updates["password"] = passwordHash
	}

	if request.CNPJ != nil {
		updates["cnpj"] = *request.CNPJ
	}

	if request.Phone != nil {
		updates["phone"] = *request.Phone
	}

	if request.Address != nil {
		updates["address"] = *request.Address
	}

	updates["updated_at"] = time.Now()

	if err := receiver.CompanyRepository.Update(request.ID, updates); err != nil {
		return nil, util.WrapError("Erro ao atualizar empresa", err, http.StatusInternalServerError)
	}

	updatedCompany, err := receiver.CompanyRepository.GetByID(request.ID)
	if err != nil {
		return nil, util.WrapError("Erro ao buscar empresa atualizada", err, http.StatusInternalServerError)
	}

	response := &contract.UpdateCompanyResponse{
		Message: "Empresa atualizada com sucesso",
		Data: contract.Company{
			ID:        updatedCompany.ID,
			Name:      updatedCompany.Name,
			Email:     updatedCompany.Email,
			CNPJ:      updatedCompany.CNPJ,
			Phone:     updatedCompany.Phone,
			Address:   updatedCompany.Address,
			CreatedAt: updatedCompany.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: updatedCompany.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	}

	return response, nil
}

// List - realiza a listagem de todas as empresas
func (receiver *CompanyService) List(filtros *model.Company) (*contract.ListCompanyResponse, error) {

	companies, err := receiver.CompanyRepository.List(filtros)
	if err != nil {
		return nil, util.WrapError("Erro ao buscar empresas", err, http.StatusInternalServerError)
	}

	var companiesResponse []contract.Company
	for _, company := range companies {
		companiesResponse = append(companiesResponse, contract.Company{
			ID:        company.ID,
			Name:      company.Name,
			Email:     company.Email,
			CNPJ:      company.CNPJ,
			Phone:     company.Phone,
			Address:   company.Address,
			CreatedAt: company.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: company.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	response := &contract.ListCompanyResponse{
		Message: "Empresas listadas com sucesso",
		Data:    companiesResponse,
	}

	return response, nil
}

// Get - realiza a busca de uma empresa por ID
func (receiver *CompanyService) Get(ID int) (*contract.GetCompanyResponse, error) {

	company, err := receiver.CompanyRepository.GetByID(ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, util.WrapError("Empresa não encontrada", nil, http.StatusNotFound)
		}
		return nil, util.WrapError("Erro ao buscar empresa", err, http.StatusInternalServerError)
	}

	response := &contract.GetCompanyResponse{
		Message: "Empresa obtida com sucesso",
		Data: contract.Company{
			ID:        company.ID,
			Name:      company.Name,
			Email:     company.Email,
			CNPJ:      company.CNPJ,
			Phone:     company.Phone,
			Address:   company.Address,
			CreatedAt: company.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: company.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	}

	return response, nil
}
