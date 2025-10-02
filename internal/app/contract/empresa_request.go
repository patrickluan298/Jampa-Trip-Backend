package contract

import (
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/jampa_trip/internal/pkg/util"
)

// LoginEmpresaRequest - objeto de request do enpoint de login de empresa
type LoginEmpresaRequest struct {
	Email string `json:"email"`
	Senha string `json:"senha"`
}

// Validate - valida os campos da requisição
func (receiver LoginEmpresaRequest) Validate() error {

	err := validation.ValidateStruct(&receiver,
		validation.Field(&receiver.Email, validation.Required, validation.Match(util.COD_03), validation.Length(1, 40)),
		validation.Field(&receiver.Senha, validation.Required, validation.Match(util.COD_07), validation.Length(8, 50)),
	)

	if err != nil {
		return util.WrapError(util.FormatarErroValidacao(err).Error(), err, http.StatusUnprocessableEntity)
	}

	return nil
}

// CadastrarEmpresaRequest - objeto de request do endpoint de cadastro de empresa
type CadastrarEmpresaRequest struct {
	Nome           string `json:"nome"`
	Email          string `json:"email"`
	Senha          string `json:"senha"`
	ConfirmarSenha string `json:"confirmar_senha"`
	CNPJ           string `json:"cnpj"`
	Telefone       string `json:"telefone"`
	Endereco       string `json:"endereco"`
}

// Validate - valida os campos da requisição de cadastro
func (receiver CadastrarEmpresaRequest) Validate() error {
	err := validation.ValidateStruct(&receiver,
		validation.Field(&receiver.Nome, validation.Required, validation.Length(2, 100)),
		validation.Field(&receiver.Email, validation.Required, validation.Match(util.COD_03), validation.Length(1, 40)),
		validation.Field(&receiver.Senha, validation.Required, validation.Match(util.COD_07), validation.Length(8, 50)),
		validation.Field(&receiver.ConfirmarSenha, validation.Required, validation.Match(util.COD_07), validation.Length(8, 50)),
		validation.Field(&receiver.CNPJ, validation.Required, validation.Match(util.COD_12)),
		validation.Field(&receiver.Telefone, validation.Required, validation.Match(util.COD_11)),
		validation.Field(&receiver.Endereco, validation.Required, validation.Length(10, 100)),
	)

	if err != nil {
		return util.WrapError(util.FormatarErroValidacao(err).Error(), err, http.StatusUnprocessableEntity)
	}

	if err := util.ValidaCNPJ(receiver.CNPJ); err != nil {
		return util.WrapError(err.Error(), err, http.StatusUnprocessableEntity)
	}

	if err := util.ValidaSegurancaSenha(receiver.Senha); err != nil {
		return util.WrapError(err.Error(), err, http.StatusUnprocessableEntity)
	}

	return nil
}

// AtualizarEmpresaRequest - objeto de request do endpoint de atualização de empresa
type AtualizarEmpresaRequest struct {
	ID       int    `json:"id"`
	Nome     string `json:"nome"`
	Email    string `json:"email"`
	Senha    string `json:"senha"`
	CNPJ     string `json:"cnpj"`
	Telefone string `json:"telefone"`
	Endereco string `json:"endereco"`
}

// Validate - valida os campos da requisição de atualização
func (receiver AtualizarEmpresaRequest) Validate() error {
	err := validation.ValidateStruct(&receiver,
		validation.Field(&receiver.ID, validation.Required, validation.Min(1)),
		validation.Field(&receiver.Nome, validation.Required, validation.Length(2, 100)),
		validation.Field(&receiver.Email, validation.Required, validation.Match(util.COD_03), validation.Length(1, 40)),
		validation.Field(&receiver.Senha, validation.Required, validation.Match(util.COD_07), validation.Length(8, 50)),
		validation.Field(&receiver.CNPJ, validation.Required, validation.Match(util.COD_12)),
		validation.Field(&receiver.Telefone, validation.Required, validation.Match(util.COD_11)),
		validation.Field(&receiver.Endereco, validation.Required, validation.Length(10, 100)),
	)

	if err != nil {
		return util.WrapError(util.FormatarErroValidacao(err).Error(), err, http.StatusUnprocessableEntity)
	}

	if err := util.ValidaCNPJ(receiver.CNPJ); err != nil {
		return util.WrapError(err.Error(), err, http.StatusUnprocessableEntity)
	}

	if err := util.ValidaSegurancaSenha(receiver.Senha); err != nil {
		return util.WrapError(err.Error(), err, http.StatusUnprocessableEntity)
	}

	return nil
}
