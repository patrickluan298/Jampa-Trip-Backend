package contract

import "time"

// PagamentoResponse - representa a resposta de um pagamento
type PagamentoResponse struct {
	ID                   int     `json:"id"`
	ClienteID            int     `json:"cliente_id"`
	EmpresaID            int     `json:"empresa_id"`
	MercadoPagoOrderID   string  `json:"mercado_pago_order_id"`
	MercadoPagoPaymentID string  `json:"mercado_pago_payment_id"`
	Status               string  `json:"status"`
	StatusDetail         string  `json:"status_detail"`
	Valor                float64 `json:"valor"`
	Moeda                string  `json:"moeda"`
	MetodoPagamento      string  `json:"metodo_pagamento"`
	Descricao            string  `json:"descricao"`
	NumeroParcelas       int     `json:"numero_parcelas"`

	// Campos específicos de cartão
	LastFourDigits            string  `json:"last_four_digits,omitempty"`
	FirstSixDigits            string  `json:"first_six_digits,omitempty"`
	PaymentMethodID           string  `json:"payment_method_id,omitempty"`
	IssuerID                  string  `json:"issuer_id,omitempty"`
	CardholderName            string  `json:"cardholder_name,omitempty"`
	Captured                  bool    `json:"captured"`
	TransactionAmountRefunded float64 `json:"transaction_amount_refunded"`

	TokenCartao string `json:"token_cartao,omitempty"`
	ChavePIX    string `json:"chave_pix,omitempty"`
	QRCode      string `json:"qr_code,omitempty"`

	MomentoCriacao      time.Time  `json:"momento_criacao"`
	MomentoAtualizacao  time.Time  `json:"momento_atualizacao"`
	MomentoAprovacao    *time.Time `json:"momento_aprovacao,omitempty"`
	MomentoCancelamento *time.Time `json:"momento_cancelamento,omitempty"`
	MomentoAutorizacao  *time.Time `json:"momento_autorizacao,omitempty"`
	MomentoCaptura      *time.Time `json:"momento_captura,omitempty"`

	StatusDisplay          string `json:"status_display"`
	MetodoPagamentoDisplay string `json:"metodo_pagamento_display"`
}

// AutorizarCartaoResponse - representa a resposta da autorização de pagamento com cartão
type AutorizarCartaoResponse struct {
	Pagamento PagamentoResponse `json:"pagamento"`
	Message   string            `json:"message"`
}

// CapturarPagamentoResponse - representa a resposta da captura de pagamento
type CapturarPagamentoResponse struct {
	Pagamento PagamentoResponse `json:"pagamento"`
	Message   string            `json:"message"`
}

// CancelarPagamentoResponse - representa a resposta do cancelamento de pagamento
type CancelarPagamentoResponse struct {
	Pagamento PagamentoResponse `json:"pagamento"`
	Message   string            `json:"message"`
}

// ReembolsarPagamentoResponse - representa a resposta do reembolso de pagamento
type ReembolsarPagamentoResponse struct {
	Pagamento      PagamentoResponse `json:"pagamento"`
	RefundID       int64             `json:"refund_id"`
	RefundedAmount float64           `json:"refunded_amount"`
	Message        string            `json:"message"`
}

// PaymentMethodsResponse - resposta com lista de meios de pagamento
type PaymentMethodsResponse struct {
	Methods []PaymentMethodInfo `json:"methods"`
}

// PaymentMethodInfo - informações de um meio de pagamento
type PaymentMethodInfo struct {
	ID                string  `json:"id"`
	Name              string  `json:"name"`
	PaymentTypeID     string  `json:"payment_type_id"`
	Status            string  `json:"status"`
	MinAllowedAmount  float64 `json:"min_allowed_amount,omitempty"`
	MaxAllowedAmount  float64 `json:"max_allowed_amount,omitempty"`
	AccreditationTime int     `json:"accreditation_time,omitempty"`
}
