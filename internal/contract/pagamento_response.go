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

// CriarPagamentoCartaoCreditoResponse - representa a resposta da criação de pagamento com cartão de crédito
type CriarPagamentoCartaoCreditoResponse struct {
	Pagamento PagamentoResponse `json:"pagamento"`
	Message   string            `json:"message"`
}

// CriarPagamentoCartaoDebitoResponse - representa a resposta da criação de pagamento com cartão de débito
type CriarPagamentoCartaoDebitoResponse struct {
	Pagamento PagamentoResponse `json:"pagamento"`
	Message   string            `json:"message"`
}

// CriarPagamentoPIXResponse - representa a resposta da criação de pagamento com PIX
type CriarPagamentoPIXResponse struct {
	Pagamento    PagamentoResponse `json:"pagamento"`
	Message      string            `json:"message"`
	QRCode       string            `json:"qr_code,omitempty"`
	QRCodeBase64 string            `json:"qr_code_base64,omitempty"`
	TicketURL    string            `json:"ticket_url,omitempty"`
}

// BuscarPagamentosResponse - representa a resposta da busca de pagamentos
type BuscarPagamentosResponse struct {
	Pagamentos []PagamentoResponse `json:"pagamentos"`
	Total      int                 `json:"total"`
	Offset     int                 `json:"offset"`
	Limit      int                 `json:"limit"`
	HasMore    bool                `json:"has_more"`
}

// ObterPagamentoResponse - representa a resposta de obter pagamento por ID
type ObterPagamentoResponse struct {
	Pagamento PagamentoResponse `json:"pagamento"`
}

// AtualizarPagamentoResponse - representa a resposta de atualizar pagamento
type AtualizarPagamentoResponse struct {
	Pagamento PagamentoResponse `json:"pagamento"`
	Message   string            `json:"message"`
}
