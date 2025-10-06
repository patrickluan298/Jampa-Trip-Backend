package contract

import (
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/jampa_trip/pkg/util"
)

// CreateCompanyRequest - objeto de request do endpoint de cadastro de empresa
type CreateCompanyRequest struct {
	Name            string `json:"name"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
	CNPJ            string `json:"cnpj"`
	Phone           string `json:"phone"`
	Address         string `json:"address"`
}

// Validate - valida os campos da requisição de cadastro
func (receiver CreateCompanyRequest) Validate() error {
	err := validation.ValidateStruct(&receiver,
		validation.Field(&receiver.Name, validation.Required, validation.Length(2, 100)),
		validation.Field(&receiver.Email, validation.Required, validation.Match(util.COD_03), validation.Length(1, 40)),
		validation.Field(&receiver.Password, validation.Required, validation.Match(util.COD_07), validation.Length(8, 50)),
		validation.Field(&receiver.ConfirmPassword, validation.Required, validation.Match(util.COD_07), validation.Length(8, 50)),
		validation.Field(&receiver.CNPJ, validation.Required, validation.Match(util.COD_12)),
		validation.Field(&receiver.Phone, validation.Required, validation.Match(util.COD_11)),
		validation.Field(&receiver.Address, validation.Required, validation.Length(10, 100)),
	)

	if err != nil {
		return util.WrapError(util.FormatarErroValidacao(err).Error(), err, http.StatusUnprocessableEntity)
	}

	if err := util.ValidaCNPJ(receiver.CNPJ); err != nil {
		return util.WrapError(err.Error(), err, http.StatusUnprocessableEntity)
	}

	if err := util.ValidaSegurancaSenha(receiver.Password); err != nil {
		return util.WrapError(err.Error(), err, http.StatusUnprocessableEntity)
	}

	return nil
}

// UpdateCompanyRequest - objeto de request do endpoint de atualização de empresa
type UpdateCompanyRequest struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Password string `json:"password"`
	CNPJ    string `json:"cnpj"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
}

// Validate - valida os campos da requisição de atualização
func (receiver UpdateCompanyRequest) Validate() error {
	err := validation.ValidateStruct(&receiver,
		validation.Field(&receiver.ID, validation.Required, validation.Min(1)),
		validation.Field(&receiver.Name, validation.Required, validation.Length(2, 100)),
		validation.Field(&receiver.Email, validation.Required, validation.Match(util.COD_03), validation.Length(1, 40)),
		validation.Field(&receiver.Password, validation.Required, validation.Match(util.COD_07), validation.Length(8, 50)),
		validation.Field(&receiver.CNPJ, validation.Required, validation.Match(util.COD_12)),
		validation.Field(&receiver.Phone, validation.Required, validation.Match(util.COD_11)),
		validation.Field(&receiver.Address, validation.Required, validation.Length(10, 100)),
	)

	if err != nil {
		return util.WrapError(util.FormatarErroValidacao(err).Error(), err, http.StatusUnprocessableEntity)
	}

	if err := util.ValidaCNPJ(receiver.CNPJ); err != nil {
		return util.WrapError(err.Error(), err, http.StatusUnprocessableEntity)
	}

	if err := util.ValidaSegurancaSenha(receiver.Password); err != nil {
		return util.WrapError(err.Error(), err, http.StatusUnprocessableEntity)
	}

	return nil
}