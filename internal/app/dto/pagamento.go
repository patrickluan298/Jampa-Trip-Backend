package dto

import (
	"github.com/jampa_trip/internal/app/model"
	"github.com/jampa_trip/internal/app/types"
)

// PaymentDisplayDTO - DTO para apresentação de pagamentos
type PaymentDisplayDTO struct {
	*model.Pagamento
}

// GetStatusDisplay - retorna uma representação amigável do status
func (d *PaymentDisplayDTO) GetStatusDisplay() string {
	switch d.Status {
	case string(types.StatusPending):
		return "Pendente"
	case string(types.StatusApproved):
		return "Aprovado"
	case string(types.StatusAuthorized):
		return "Autorizado"
	case string(types.StatusInProcess):
		return "Em Processamento"
	case string(types.StatusInMediation):
		return "Em Mediação"
	case string(types.StatusRejected):
		return "Rejeitado"
	case string(types.StatusCancelled):
		return "Cancelado"
	case string(types.StatusRefunded):
		return "Reembolsado"
	case string(types.StatusChargedBack):
		return "Estornado"
	default:
		return "Desconhecido"
	}
}

// GetPaymentMethodDisplay - retorna uma representação amigável do método de pagamento
func (d *PaymentDisplayDTO) GetPaymentMethodDisplay() string {
	switch d.MetodoPagamento {
	case string(types.MetodoCartaoCredito):
		return "Cartão de Crédito"
	case string(types.MetodoCartaoDebito):
		return "Cartão de Débito"
	case string(types.MetodoPIX):
		return "PIX"
	case string(types.MetodoBoleto):
		return "Boleto"
	case string(types.MetodoPec):
		return "PEC"
	default:
		return "Desconhecido"
	}
}
