package service

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/jampa_trip/internal/app/contract"
	"github.com/jampa_trip/internal/app/repository"
	"github.com/jampa_trip/internal/pkg/util"
	"golang.org/x/crypto/bcrypt"
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

	fornecedor, err := receiver.FornecedorRepository.GetFornecedor(request.Email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, util.WrapError("Email e/ou senha incorretos", nil, http.StatusUnauthorized)
		}
		return nil, util.WrapError("Erro ao buscar usuário", err, http.StatusInternalServerError)
	}

	if !receiver.verifyPassword(request.Senha, fornecedor.Senha) {
		return nil, util.WrapError("Email e/ou senha incorretos", nil, http.StatusUnauthorized)
	}

	token, err := receiver.generateToken()
	if err != nil {
		return nil, util.WrapError("Erro ao gerar token", err, http.StatusInternalServerError)
	}

	response := &contract.LoginFornecedorResponse{
		StatusCode: http.StatusOK,
		Message:    "Login realizado com sucesso",
		Token:      token,
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

// verifyPassword verifica se a senha fornecida corresponde ao hash armazenado
func (receiver *FornecedorService) verifyPassword(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// generateToken gera um token aleatório para autenticação
func (receiver *FornecedorService) generateToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// HashPassword gera um hash da senha usando bcrypt
func (receiver *FornecedorService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
