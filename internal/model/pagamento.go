package model

import "time"

// Pagamento - representa a entidade de pagamento
type Pagamento struct {
	ID                   int     `gorm:"column:id;primaryKey;autoIncrement"`
	ClienteID            int     `gorm:"column:cliente_id;not null"`
	EmpresaID            int     `gorm:"column:empresa_id;not null"`
	MercadoPagoOrderID   string  `gorm:"column:mercado_pago_order_id;uniqueIndex"`
	MercadoPagoPaymentID string  `gorm:"column:mercado_pago_payment_id;index"`
	Status               string  `gorm:"column:status;not null;default:'pending'"`
	StatusDetail         string  `gorm:"column:status_detail"`
	Valor                float64 `gorm:"column:valor;not null;type:decimal(10,2)"`
	Moeda                string  `gorm:"column:moeda;not null;default:'BRL'"`
	MetodoPagamento      string  `gorm:"column:metodo_pagamento;not null"`
	Descricao            string  `gorm:"column:descricao"`
	NumeroParcelas       int     `gorm:"column:numero_parcelas;default:1"`
	TokenCartao          string  `gorm:"column:token_cartao"`
	ChavePIX             string  `gorm:"column:chave_pix"`
	QRCode               string  `gorm:"column:qr_code;type:text"`

	// Campos específicos para cartão de crédito (dados não sensíveis)
	LastFourDigits            string  `gorm:"column:last_four_digits"`
	FirstSixDigits            string  `gorm:"column:first_six_digits"`
	PaymentMethodID           string  `gorm:"column:payment_method_id"`
	IssuerID                  string  `gorm:"column:issuer_id"`
	CardholderName            string  `gorm:"column:cardholder_name"`
	Captured                  bool    `gorm:"column:captured;default:false"`
	TransactionAmountRefunded float64 `gorm:"column:transaction_amount_refunded;type:decimal(10,2);default:0"`

	MomentoCriacao      time.Time  `gorm:"column:momento_criacao;not null;default:CURRENT_TIMESTAMP"`
	MomentoAtualizacao  time.Time  `gorm:"column:momento_atualizacao;not null;default:CURRENT_TIMESTAMP"`
	MomentoAprovacao    *time.Time `gorm:"column:momento_aprovacao"`
	MomentoCancelamento *time.Time `gorm:"column:momento_cancelamento"`
	MomentoAutorizacao  *time.Time `gorm:"column:momento_autorizacao"`
	MomentoCaptura      *time.Time `gorm:"column:momento_captura"`

	// Relacionamentos
	Cliente Cliente `gorm:"foreignKey:ClienteID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Empresa Empresa `gorm:"foreignKey:EmpresaID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

// TableName - especifica o nome da tabela no banco de dados
func (Pagamento) TableName() string {
	return "pagamentos"
}

// StatusPagamento - define os possíveis status de um pagamento
type StatusPagamento string

const (
	StatusPending     StatusPagamento = "pending"
	StatusApproved    StatusPagamento = "approved"
	StatusAuthorized  StatusPagamento = "authorized"
	StatusInProcess   StatusPagamento = "in_process"
	StatusInMediation StatusPagamento = "in_mediation"
	StatusRejected    StatusPagamento = "rejected"
	StatusCancelled   StatusPagamento = "cancelled"
	StatusRefunded    StatusPagamento = "refunded"
	StatusChargedBack StatusPagamento = "charged_back"
)

// MetodoPagamento - define os métodos de pagamento suportados
type MetodoPagamento string

const (
	MetodoCartaoCredito MetodoPagamento = "credit_card"
	MetodoCartaoDebito  MetodoPagamento = "debit_card"
	MetodoPIX           MetodoPagamento = "pix"
)

// Moeda - define as moedas suportadas
type Moeda string

const (
	MoedaBRL Moeda = "BRL"
	MoedaUSD Moeda = "USD"
	MoedaEUR Moeda = "EUR"
	MoedaARS Moeda = "ARS"
	MoedaCLP Moeda = "CLP"
	MoedaCOP Moeda = "COP"
	MoedaMXN Moeda = "MXN"
	MoedaPEN Moeda = "PEN"
	MoedaUYU Moeda = "UYU"
)

// Métodos de validação para Pagamento
func (p *Pagamento) IsValid() bool {
	return p.Valor > 0 && p.Status != "" && p.MetodoPagamento != ""
}

func (p *Pagamento) IsApproved() bool {
	return p.Status == string(StatusApproved)
}

func (p *Pagamento) IsPending() bool {
	return p.Status == string(StatusPending)
}

func (p *Pagamento) IsCancelled() bool {
	return p.Status == string(StatusCancelled)
}

func (p *Pagamento) IsRejected() bool {
	return p.Status == string(StatusRejected)
}

func (p *Pagamento) IsAuthorized() bool {
	return p.Status == string(StatusAuthorized)
}

func (p *Pagamento) IsCaptured() bool {
	return p.Captured && p.Status == string(StatusApproved)
}

func (p *Pagamento) UpdateStatus(status StatusPagamento) {
	p.Status = string(status)
	p.MomentoAtualizacao = time.Now()
}

func (p *Pagamento) GetStatusDisplay() string {
	switch p.Status {
	case string(StatusPending):
		return "Pendente"
	case string(StatusApproved):
		return "Aprovado"
	case string(StatusAuthorized):
		return "Autorizado"
	case string(StatusInProcess):
		return "Em Processamento"
	case string(StatusInMediation):
		return "Em Mediação"
	case string(StatusRejected):
		return "Rejeitado"
	case string(StatusCancelled):
		return "Cancelado"
	case string(StatusRefunded):
		return "Reembolsado"
	case string(StatusChargedBack):
		return "Estornado"
	default:
		return "Desconhecido"
	}
}

func (p *Pagamento) GetMetodoPagamentoDisplay() string {
	switch p.MetodoPagamento {
	case string(MetodoCartaoCredito):
		return "Cartão de Crédito"
	case string(MetodoCartaoDebito):
		return "Cartão de Débito"
	case string(MetodoPIX):
		return "PIX"
	default:
		return "Desconhecido"
	}
}

// Funções de validação
func IsValidStatus(status StatusPagamento) bool {
	switch status {
	case StatusPending, StatusApproved, StatusAuthorized, StatusInProcess,
		StatusInMediation, StatusRejected, StatusCancelled, StatusRefunded, StatusChargedBack:
		return true
	default:
		return false
	}
}

func IsValidPaymentMethod(method MetodoPagamento) bool {
	switch method {
	case MetodoCartaoCredito, MetodoCartaoDebito, MetodoPIX:
		return true
	default:
		return false
	}
}

func IsValidCurrency(currency Moeda) bool {
	switch currency {
	case MoedaBRL, MoedaUSD, MoedaEUR, MoedaARS, MoedaCLP, MoedaCOP, MoedaMXN, MoedaPEN, MoedaUYU:
		return true
	default:
		return false
	}
}
