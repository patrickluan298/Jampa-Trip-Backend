package contract

import (
	"net/http"
	"strings"

	"github.com/jampa_trip/pkg/util"
)

// CreateCartaoRequest - representa a requisição para criar um cartão
type CreateCartaoRequest struct {
	Token           string            `json:"token" validate:"required"`
	PaymentMethodID string            `json:"payment_method_id" validate:"required"`
	IssuerID        string            `json:"issuer_id,omitempty"`
	Cardholder      CardholderRequest `json:"cardholder" validate:"required"`
	Metadata        map[string]string `json:"metadata,omitempty"`
}

// CardholderRequest - representa o portador do cartão
type CardholderRequest struct {
	Name           string                      `json:"name" validate:"required"`
	Identification CartaoIdentificationRequest `json:"identification" validate:"required"`
}

// CartaoIdentificationRequest - representa a identificação do portador
type CartaoIdentificationRequest struct {
	Type   string `json:"type" validate:"required"`
	Number string `json:"number" validate:"required"`
}

// UpdateCartaoRequest - representa a requisição para atualizar um cartão
type UpdateCartaoRequest struct {
	Cardholder CardholderRequest `json:"cardholder,omitempty"`
	Metadata   map[string]string `json:"metadata,omitempty"`
	Default    bool              `json:"default,omitempty"`
}

// Validate - valida os dados da requisição de criação de cartão
func (r *CreateCartaoRequest) Validate() error {
	if strings.TrimSpace(r.Token) == "" {
		return util.WrapError("token é obrigatório", nil, http.StatusUnprocessableEntity)
	}
	if strings.TrimSpace(r.PaymentMethodID) == "" {
		return util.WrapError("payment_method_id é obrigatório", nil, http.StatusUnprocessableEntity)
	}
	if strings.TrimSpace(r.Cardholder.Name) == "" {
		return util.WrapError("nome do portador é obrigatório", nil, http.StatusUnprocessableEntity)
	}
	if strings.TrimSpace(r.Cardholder.Identification.Type) == "" {
		return util.WrapError("tipo de identificação é obrigatório", nil, http.StatusUnprocessableEntity)
	}
	if strings.TrimSpace(r.Cardholder.Identification.Number) == "" {
		return util.WrapError("número de identificação é obrigatório", nil, http.StatusUnprocessableEntity)
	}
	return nil
}

// Validate - valida os dados da requisição de atualização de cartão
func (r *UpdateCartaoRequest) Validate() error {

	if r.Cardholder.Name == "" && r.Cardholder.Identification.Type == "" &&
		r.Cardholder.Identification.Number == "" && len(r.Metadata) == 0 {
		return util.WrapError("pelo menos um campo deve ser fornecido para atualização", nil, http.StatusUnprocessableEntity)
	}

	if r.Cardholder.Identification.Type != "" || r.Cardholder.Identification.Number != "" {
		if strings.TrimSpace(r.Cardholder.Identification.Type) == "" {
			return util.WrapError("tipo de identificação é obrigatório quando número é fornecido", nil, http.StatusUnprocessableEntity)
		}
		if strings.TrimSpace(r.Cardholder.Identification.Number) == "" {
			return util.WrapError("número de identificação é obrigatório quando tipo é fornecido", nil, http.StatusUnprocessableEntity)
		}
	}

	return nil
}
