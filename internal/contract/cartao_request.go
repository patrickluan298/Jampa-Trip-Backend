package contract

import (
	"errors"
	"strings"
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
		return errors.New("token é obrigatório")
	}
	if strings.TrimSpace(r.PaymentMethodID) == "" {
		return errors.New("payment_method_id é obrigatório")
	}
	if strings.TrimSpace(r.Cardholder.Name) == "" {
		return errors.New("nome do portador é obrigatório")
	}
	if strings.TrimSpace(r.Cardholder.Identification.Type) == "" {
		return errors.New("tipo de identificação é obrigatório")
	}
	if strings.TrimSpace(r.Cardholder.Identification.Number) == "" {
		return errors.New("número de identificação é obrigatório")
	}
	return nil
}

// Validate - valida os dados da requisição de atualização de cartão
func (r *UpdateCartaoRequest) Validate() error {

	if r.Cardholder.Name == "" && r.Cardholder.Identification.Type == "" &&
		r.Cardholder.Identification.Number == "" && len(r.Metadata) == 0 {
		return errors.New("pelo menos um campo deve ser fornecido para atualização")
	}

	if r.Cardholder.Identification.Type != "" || r.Cardholder.Identification.Number != "" {
		if strings.TrimSpace(r.Cardholder.Identification.Type) == "" {
			return errors.New("tipo de identificação é obrigatório quando número é fornecido")
		}
		if strings.TrimSpace(r.Cardholder.Identification.Number) == "" {
			return errors.New("número de identificação é obrigatório quando tipo é fornecido")
		}
	}

	return nil
}
