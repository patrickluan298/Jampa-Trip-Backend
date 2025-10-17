package contract

import "time"

// PaymentResponse - representa a resposta de um pagamento
type PaymentResponse struct {
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

// CreateCreditCardPaymentResponse - representa a resposta da criação de pagamento com cartão de crédito
type CreateCreditCardPaymentResponse struct {
	Pagamento PaymentResponse `json:"pagamento"`
	Message   string          `json:"message"`
}

// CreateDebitCardPaymentResponse - representa a resposta da criação de pagamento com cartão de débito
type CreateDebitCardPaymentResponse struct {
	Pagamento PaymentResponse `json:"pagamento"`
	Message   string          `json:"message"`
}

// CreatePIXPaymentResponse - representa a resposta da criação de pagamento com PIX
type CreatePIXPaymentResponse struct {
	Pagamento    PaymentResponse `json:"pagamento"`
	Message      string          `json:"message"`
	QRCode       string          `json:"qr_code,omitempty"`
	QRCodeBase64 string          `json:"qr_code_base64,omitempty"`
	TicketURL    string          `json:"ticket_url,omitempty"`
}

// ListPaymentsResponse - representa a resposta da busca de pagamentos
type ListPaymentsResponse struct {
	Pagamentos []PaymentResponse `json:"pagamentos"`
	Total      int               `json:"total"`
	Offset     int               `json:"offset"`
	Limit      int               `json:"limit"`
	HasMore    bool              `json:"has_more"`
}

// GetPaymentResponse - representa a resposta de obter pagamento por ID
type GetPaymentResponse struct {
	Pagamento PaymentResponse `json:"pagamento"`
}

// UpdatePaymentResponse - representa a resposta de atualizar pagamento
type UpdatePaymentResponse struct {
	Pagamento PaymentResponse `json:"pagamento"`
	Message   string          `json:"message"`
}
