package service

import (
	"net/http"

	"github.com/jampa_trip/internal/contract"
	"github.com/jampa_trip/internal/repository"
	"github.com/jampa_trip/pkg/auth"
	"github.com/jampa_trip/pkg/util"
	"gorm.io/gorm"
)

// LoginService - objeto de contexto para login
type LoginService struct {
	CompanyRepository *repository.CompanyRepository
	ClientRepository  *repository.ClientRepository
}

// LoginServiceNew - construtor do objeto
func LoginServiceNew(DB *gorm.DB) *LoginService {
	return &LoginService{
		CompanyRepository: repository.CompanyRepositoryNew(DB),
		ClientRepository:  repository.ClientRepositoryNew(DB),
	}
}

// Login - realiza a autenticação
func (receiver *LoginService) Login(request *contract.LoginRequest) (*contract.LoginResponse, error) {

	company, err := receiver.CompanyRepository.GetByEmail(request.Email)
	if err == nil {
		if util.VerificaSenha(request.Password, company.Password) {
			tokenPair, err := auth.GenerateTokenPair(company.ID, "company", company.Email)
			if err != nil {
				return nil, util.WrapError("erro ao gerar tokens JWT", err, http.StatusInternalServerError)
			}

			tokenStore := auth.NewRedisTokenStore()
			err = tokenStore.StoreTokenPair(company.ID, "company", tokenPair.AccessToken, tokenPair.RefreshToken)
			if err != nil {
				return nil, util.WrapError("erro ao armazenar tokens no Redis", err, http.StatusInternalServerError)
			}

			response := &contract.LoginResponse{
				Message: "Login realizado com sucesso",
				Type:    "company",
				Data: contract.UserLoginData{
					ID:    company.ID,
					Name:  company.Name,
					Email: company.Email,
				},
				AccessToken:  tokenPair.AccessToken,
				RefreshToken: tokenPair.RefreshToken,
				ExpiresIn:    tokenPair.ExpiresIn,
			}
			return response, nil
		}
		return nil, util.WrapError("Email e/ou senha incorretos", nil, http.StatusUnauthorized)
	}

	client, err := receiver.ClientRepository.GetByEmail(request.Email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, util.WrapError("Email e/ou senha incorretos", nil, http.StatusUnauthorized)
		}
		return nil, util.WrapError("Erro ao buscar usuário", err, http.StatusInternalServerError)
	}

	if !util.VerificaSenha(request.Password, client.Password) {
		return nil, util.WrapError("Email e/ou senha incorretos", nil, http.StatusUnauthorized)
	}

	tokenPair, err := auth.GenerateTokenPair(client.ID, "client", client.Email)
	if err != nil {
		return nil, util.WrapError("erro ao gerar tokens JWT", err, http.StatusInternalServerError)
	}

	tokenStore := auth.NewRedisTokenStore()
	err = tokenStore.StoreTokenPair(client.ID, "client", tokenPair.AccessToken, tokenPair.RefreshToken)
	if err != nil {
		return nil, util.WrapError("erro ao armazenar tokens no Redis", err, http.StatusInternalServerError)
	}

	response := &contract.LoginResponse{
		Message: "Login realizado com sucesso",
		Type:    "client",
		Data: contract.UserLoginData{
			ID:    client.ID,
			Name:  client.Name,
			Email: client.Email,
		},
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    tokenPair.ExpiresIn,
	}

	return response, nil
}
