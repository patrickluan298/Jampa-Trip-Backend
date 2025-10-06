package contract

import (
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/jampa_trip/pkg/util"
)

// CreateClientRequest - objeto de request do endpoint de cadastro de cliente
type CreateClientRequest struct {
	Name            string `json:"name"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
	CPF             string `json:"cpf"`
	Phone           string `json:"phone"`
	BirthDate       string `json:"birth_date"`
}

// Validate - valida os campos da requisição de cadastro
func (receiver CreateClientRequest) Validate() error {
	err := validation.ValidateStruct(&receiver,
		validation.Field(&receiver.Name, validation.Required, validation.Length(2, 100)),
		validation.Field(&receiver.Email, validation.Required, validation.Match(util.COD_03), validation.Length(1, 40)),
		validation.Field(&receiver.Password, validation.Required, validation.Match(util.COD_07), validation.Length(8, 50)),
		validation.Field(&receiver.ConfirmPassword, validation.Required, validation.Match(util.COD_07), validation.Length(8, 50)),
		validation.Field(&receiver.CPF, validation.Required, validation.Match(util.COD_04)),
		validation.Field(&receiver.Phone, validation.Required, validation.Match(util.COD_11)),
		validation.Field(&receiver.BirthDate, validation.Required, validation.Match(util.COD_06)),
	)

	if err != nil {
		return util.WrapError(util.FormatarErroValidacao(err).Error(), err, http.StatusUnprocessableEntity)
	}

	if err := util.ValidaCPF(receiver.CPF); err != nil {
		return util.WrapError(err.Error(), err, http.StatusUnprocessableEntity)
	}

	if err := util.ValidaSegurancaSenha(receiver.Password); err != nil {
		return util.WrapError(err.Error(), err, http.StatusUnprocessableEntity)
	}

	return nil
}

// UpdateClientRequest - objeto de request do endpoint de atualização de cliente
type UpdateClientRequest struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	CPF       string `json:"cpf"`
	Phone     string `json:"phone"`
	BirthDate string `json:"birth_date"`
}

// Validate - valida os campos da requisição de atualização
func (receiver UpdateClientRequest) Validate() error {
	err := validation.ValidateStruct(&receiver,
		validation.Field(&receiver.ID, validation.Required, validation.Min(1)),
		validation.Field(&receiver.Name, validation.Required, validation.Length(2, 100)),
		validation.Field(&receiver.Email, validation.Required, validation.Match(util.COD_03), validation.Length(1, 40)),
		validation.Field(&receiver.Password, validation.Required, validation.Match(util.COD_07), validation.Length(8, 50)),
		validation.Field(&receiver.CPF, validation.Required, validation.Match(util.COD_04)),
		validation.Field(&receiver.Phone, validation.Required, validation.Match(util.COD_11)),
		validation.Field(&receiver.BirthDate, validation.Required, validation.Match(util.COD_06)),
	)

	if err != nil {
		return util.WrapError(util.FormatarErroValidacao(err).Error(), err, http.StatusUnprocessableEntity)
	}

	if err := util.ValidaCPF(receiver.CPF); err != nil {
		return util.WrapError(err.Error(), err, http.StatusUnprocessableEntity)
	}

	if err := util.ValidaSegurancaSenha(receiver.Password); err != nil {
		return util.WrapError(err.Error(), err, http.StatusUnprocessableEntity)
	}

	return nil
}
