package contract

import (
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/jampa_trip/pkg/util"
)

// LoginRequest - objeto de request do endpoint de login unificado
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Validate - valida os campos da requisição
func (receiver LoginRequest) Validate() error {

	err := validation.ValidateStruct(&receiver,
		validation.Field(&receiver.Email, validation.Required, validation.Match(util.COD_03), validation.Length(1, 40)),
		validation.Field(&receiver.Password, validation.Required, validation.Match(util.COD_07), validation.Length(8, 50)),
	)

	if err != nil {
		return util.WrapError(util.FormatarErroValidacao(err).Error(), err, http.StatusUnprocessableEntity)
	}

	return nil
}
