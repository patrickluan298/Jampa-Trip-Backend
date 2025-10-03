package contract

import (
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/jampa_trip/pkg/util"
)

// LoginClienteRequest - objeto de request do endpoint de login de cliente
type LoginClienteRequest struct {
	Email string `json:"email"`
	Senha string `json:"senha"`
}

// Validate - valida os campos da requisição
func (receiver LoginClienteRequest) Validate() error {

	err := validation.ValidateStruct(&receiver,
		validation.Field(&receiver.Email, validation.Required, validation.Match(util.COD_03), validation.Length(1, 40)),
		validation.Field(&receiver.Senha, validation.Required, validation.Match(util.COD_07), validation.Length(8, 50)),
	)

	if err != nil {
		return util.WrapError(util.FormatarErroValidacao(err).Error(), err, http.StatusUnprocessableEntity)
	}

	return nil
}

// CadastrarClienteRequest - objeto de request do endpoint de cadastro de cliente
type CadastrarClienteRequest struct {
	Nome           string `json:"nome"`
	Email          string `json:"email"`
	Senha          string `json:"senha"`
	ConfirmarSenha string `json:"confirmar_senha"`
	CPF            string `json:"cpf"`
	Telefone       string `json:"telefone"`
	DataNascimento string `json:"data_nascimento"`
}

// Validate - valida os campos da requisição de cadastro
func (receiver CadastrarClienteRequest) Validate() error {
	err := validation.ValidateStruct(&receiver,
		validation.Field(&receiver.Nome, validation.Required, validation.Length(2, 100)),
		validation.Field(&receiver.Email, validation.Required, validation.Match(util.COD_03), validation.Length(1, 40)),
		validation.Field(&receiver.Senha, validation.Required, validation.Match(util.COD_07), validation.Length(8, 50)),
		validation.Field(&receiver.ConfirmarSenha, validation.Required, validation.Match(util.COD_07), validation.Length(8, 50)),
		validation.Field(&receiver.CPF, validation.Required, validation.Match(util.COD_04)),
		validation.Field(&receiver.Telefone, validation.Required, validation.Match(util.COD_11)),
		validation.Field(&receiver.DataNascimento, validation.Required, validation.Match(util.COD_06)),
	)

	if err != nil {
		return util.WrapError(util.FormatarErroValidacao(err).Error(), err, http.StatusUnprocessableEntity)
	}

	if err := util.ValidaCPF(receiver.CPF); err != nil {
		return util.WrapError(err.Error(), err, http.StatusUnprocessableEntity)
	}

	if err := util.ValidaSegurancaSenha(receiver.Senha); err != nil {
		return util.WrapError(err.Error(), err, http.StatusUnprocessableEntity)
	}

	return nil
}

// AtualizarClienteRequest - objeto de request do endpoint de atualização de cliente
type AtualizarClienteRequest struct {
	ID             int    `json:"id"`
	Nome           string `json:"nome"`
	Email          string `json:"email"`
	Senha          string `json:"senha"`
	CPF            string `json:"cpf"`
	Telefone       string `json:"telefone"`
	DataNascimento string `json:"data_nascimento"`
}

// Validate - valida os campos da requisição de atualização
func (receiver AtualizarClienteRequest) Validate() error {
	err := validation.ValidateStruct(&receiver,
		validation.Field(&receiver.ID, validation.Required, validation.Min(1)),
		validation.Field(&receiver.Nome, validation.Required, validation.Length(2, 100)),
		validation.Field(&receiver.Email, validation.Required, validation.Match(util.COD_03), validation.Length(1, 40)),
		validation.Field(&receiver.Senha, validation.Required, validation.Match(util.COD_07), validation.Length(8, 50)),
		validation.Field(&receiver.CPF, validation.Required, validation.Match(util.COD_04)),
		validation.Field(&receiver.Telefone, validation.Required, validation.Match(util.COD_11)),
		validation.Field(&receiver.DataNascimento, validation.Required, validation.Match(util.COD_06)),
	)

	if err != nil {
		return util.WrapError(util.FormatarErroValidacao(err).Error(), err, http.StatusUnprocessableEntity)
	}

	if err := util.ValidaCPF(receiver.CPF); err != nil {
		return util.WrapError(err.Error(), err, http.StatusUnprocessableEntity)
	}

	if err := util.ValidaSegurancaSenha(receiver.Senha); err != nil {
		return util.WrapError(err.Error(), err, http.StatusUnprocessableEntity)
	}

	return nil
}
