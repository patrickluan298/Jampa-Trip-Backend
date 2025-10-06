package contract

import (
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/jampa_trip/pkg/util"
)

// RefreshTokenRequest - request para renovação de token
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// Validate - valida os campos da requisição
func (receiver RefreshTokenRequest) Validate() error {
	err := validation.ValidateStruct(&receiver,
		validation.Field(&receiver.RefreshToken, validation.Required, validation.Length(1, 500)),
	)

	if err != nil {
		return util.WrapError(util.FormatarErroValidacao(err).Error(), err, http.StatusUnprocessableEntity)
	}

	return nil
}
