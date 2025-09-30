package contract

import "time"

// PagamentoResponse - representa a resposta de um pagamento
type PagamentoResponse struct {
	ID                     int        `json:"id"`
	ClienteID              int        `json:"cliente_id"`
	EmpresaID              int        `json:"empresa_id"`
	MercadoPagoOrderID     string     `json:"mercado_pago_order_id"`
	MercadoPagoPaymentID   string     `json:"mercado_pago_payment_id"`
	Status                 string     `json:"status"`
	StatusDetail           string     `json:"status_detail"`
	Valor                  float64    `json:"valor"`
	Moeda                  string     `json:"moeda"`
	MetodoPagamento        string     `json:"metodo_pagamento"`
	Descricao              string     `json:"descricao"`
	NumeroParcelas         int        `json:"numero_parcelas"`
	TokenCartao            string     `json:"token_cartao"`
	ChavePIX               string     `json:"chave_pix"`
	QRCode                 string     `json:"qr_code"`
	MomentoCriacao         time.Time  `json:"momento_criacao"`
	MomentoAtualizacao     time.Time  `json:"momento_atualizacao"`
	MomentoAprovacao       *time.Time `json:"momento_aprovacao"`
	MomentoCancelamento    *time.Time `json:"momento_cancelamento"`
	StatusDisplay          string     `json:"status_display"`
	MetodoPagamentoDisplay string     `json:"metodo_pagamento_display"`
}

// ListPagamentoResponse - representa a resposta de uma lista de pagamentos
type ListPagamentoResponse struct {
	Pagamentos []PagamentoResponse `json:"pagamentos"`
	Total      int                 `json:"total"`
	Page       int                 `json:"page"`
	Limit      int                 `json:"limit"`
	Pages      int                 `json:"pages"`
}

// CreatePagamentoResponse - representa a resposta da criação de um pagamento
type CreatePagamentoResponse struct {
	Pagamento PagamentoResponse `json:"pagamento"`
	Message   string            `json:"message"`
}

// UpdatePagamentoResponse - representa a resposta da atualização de um pagamento
type UpdatePagamentoResponse struct {
	Pagamento PagamentoResponse `json:"pagamento"`
	Message   string            `json:"message"`
}

// DeletePagamentoResponse - representa a resposta da exclusão de um pagamento
type DeletePagamentoResponse struct {
	Message string `json:"message"`
}
