package contract

import (
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/jampa_trip/internal/pkg/util"
)

// LoginFornecedorRequest - objeto de request do enpoint de login de fornecedor
type LoginFornecedorRequest struct {
	Email string `json:"email"`
	Senha string `json:"senha"`
}

// Validate valida os campos da requisição
func (receiver LoginFornecedorRequest) Validate() error {

	err := validation.ValidateStruct(&receiver,
		validation.Field(&receiver.Email, validation.Required, validation.Match(util.COD_03), validation.Length(1, 40)),
		validation.Field(&receiver.Senha, validation.Required, validation.Match(util.COD_07), validation.Length(8, 50)),
	)

	if err != nil {
		return util.WrapError(util.FormatarErroValidacao(err).Error(), err, http.StatusUnprocessableEntity)
	}

	return nil
}
