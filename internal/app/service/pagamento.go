package service

import (
	"time"

	"github.com/jampa_trip/internal/app/model"
	"github.com/jampa_trip/internal/app/types"
)

// PaymentService - serviço para lógica de negócio de pagamentos
type PaymentService struct{}

// NewPaymentService - construtor do serviço
func NewPaymentService() *PaymentService {
	return &PaymentService{}
}

// IsApproved - verifica se o pagamento foi aprovado
func (s *PaymentService) IsApproved(payment *model.Pagamento) bool {
	return payment.Status == string(types.StatusApproved)
}

// IsPending - verifica se o pagamento está pendente
func (s *PaymentService) IsPending(payment *model.Pagamento) bool {
	return payment.Status == string(types.StatusPending)
}

// CanBeCancelled - verifica se o pagamento pode ser cancelado
func (s *PaymentService) CanBeCancelled(payment *model.Pagamento) bool {
	return payment.Status == string(types.StatusPending) || payment.Status == string(types.StatusInProcess)
}

// UpdateStatus - atualiza o status do pagamento
func (s *PaymentService) UpdateStatus(payment *model.Pagamento, newStatus string) {
	payment.Status = newStatus
	payment.MomentoAtualizacao = time.Now()

	// Atualizar timestamps específicos baseado no status
	switch newStatus {
	case string(types.StatusApproved), string(types.StatusAuthorized):
		if payment.MomentoAprovacao == nil {
			now := time.Now()
			payment.MomentoAprovacao = &now
		}
	case string(types.StatusCancelled):
		if payment.MomentoCancelamento == nil {
			now := time.Now()
			payment.MomentoCancelamento = &now
		}
	}
}
