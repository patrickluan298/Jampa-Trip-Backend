package model

import "time"

// Transaction - representa uma transação associada a uma ordem
type Transaction struct {
	ID                       int       `gorm:"column:id;primaryKey;autoIncrement"`
	OrderID                  int       `gorm:"column:order_id;not null;index"`
	MercadoPagoTransactionID string    `gorm:"column:mercado_pago_transaction_id;index"`
	PaymentID                int64     `gorm:"column:payment_id;index"`
	Type                     string    `gorm:"column:type;not null"` // payment, refund
	Status                   string    `gorm:"column:status;not null;default:'pending'"`
	StatusDetail             string    `gorm:"column:status_detail"`
	Amount                   float64   `gorm:"column:amount;not null;type:decimal(10,2)"`
	Currency                 string    `gorm:"column:currency;not null;default:'BRL'"`
	PaymentMethodID          string    `gorm:"column:payment_method_id"`
	PaymentTypeID            string    `gorm:"column:payment_type_id"`
	MomentoCriacao           time.Time `gorm:"column:momento_criacao;not null;default:CURRENT_TIMESTAMP"`
	MomentoAtualizacao       time.Time `gorm:"column:momento_atualizacao;not null;default:CURRENT_TIMESTAMP"`

	// Relacionamento
	Order Order `gorm:"foreignKey:OrderID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

// TableName - especifica o nome da tabela no banco de dados
func (Transaction) TableName() string {
	return "transactions"
}

// TransactionType - define os tipos de transação
type TransactionType string

const (
	TransactionTypePayment TransactionType = "payment"
	TransactionTypeRefund  TransactionType = "refund"
)

// TransactionStatus - define os possíveis status de uma transação
type TransactionStatus string

const (
	TransactionStatusPending    TransactionStatus = "pending"
	TransactionStatusApproved   TransactionStatus = "approved"
	TransactionStatusAuthorized TransactionStatus = "authorized"
	TransactionStatusRejected   TransactionStatus = "rejected"
	TransactionStatusCancelled  TransactionStatus = "cancelled"
	TransactionStatusRefunded   TransactionStatus = "refunded"
)

// IsValid - valida se a transação possui dados válidos
func (t *Transaction) IsValid() bool {
	return t.Amount > 0 && t.Type != "" && t.Currency != ""
}

// IsPayment - verifica se é uma transação de pagamento
func (t *Transaction) IsPayment() bool {
	return t.Type == string(TransactionTypePayment)
}

// IsRefund - verifica se é uma transação de reembolso
func (t *Transaction) IsRefund() bool {
	return t.Type == string(TransactionTypeRefund)
}

// IsApproved - verifica se a transação foi aprovada
func (t *Transaction) IsApproved() bool {
	return t.Status == string(TransactionStatusApproved)
}

// IsPending - verifica se a transação está pendente
func (t *Transaction) IsPending() bool {
	return t.Status == string(TransactionStatusPending)
}

// IsRejected - verifica se a transação foi rejeitada
func (t *Transaction) IsRejected() bool {
	return t.Status == string(TransactionStatusRejected)
}

// UpdateStatus - atualiza o status da transação
func (t *Transaction) UpdateStatus(status TransactionStatus) {
	t.Status = string(status)
	t.MomentoAtualizacao = time.Now()
}

// GetStatusDisplay - retorna uma descrição amigável do status
func (t *Transaction) GetStatusDisplay() string {
	switch t.Status {
	case string(TransactionStatusPending):
		return "Pendente"
	case string(TransactionStatusApproved):
		return "Aprovado"
	case string(TransactionStatusAuthorized):
		return "Autorizado"
	case string(TransactionStatusRejected):
		return "Rejeitado"
	case string(TransactionStatusCancelled):
		return "Cancelado"
	case string(TransactionStatusRefunded):
		return "Reembolsado"
	default:
		return "Desconhecido"
	}
}

// GetTypeDisplay - retorna uma descrição amigável do tipo
func (t *Transaction) GetTypeDisplay() string {
	switch t.Type {
	case string(TransactionTypePayment):
		return "Pagamento"
	case string(TransactionTypeRefund):
		return "Reembolso"
	default:
		return "Desconhecido"
	}
}

// IsValidTransactionType - valida se o tipo é válido
func IsValidTransactionType(txType TransactionType) bool {
	switch txType {
	case TransactionTypePayment, TransactionTypeRefund:
		return true
	default:
		return false
	}
}

// IsValidTransactionStatus - valida se o status é válido
func IsValidTransactionStatus(status TransactionStatus) bool {
	switch status {
	case TransactionStatusPending, TransactionStatusApproved, TransactionStatusAuthorized,
		TransactionStatusRejected, TransactionStatusCancelled, TransactionStatusRefunded:
		return true
	default:
		return false
	}
}
