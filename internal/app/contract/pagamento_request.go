package contract

import validation "github.com/go-ozzo/ozzo-validation"

// CreatePagamentoRequest - representa a requisição para criar um pagamento
type CreatePagamentoRequest struct {
	ClienteID       int     `json:"cliente_id" validate:"required,min=1"`
	EmpresaID       int     `json:"empresa_id" validate:"required,min=1"`
	Valor           float64 `json:"valor" validate:"required,min=0.01"`
	Moeda           string  `json:"moeda" validate:"required,oneof=BRL USD EUR ARS CLP COP MXN PEN UYU"`
	MetodoPagamento string  `json:"metodo_pagamento" validate:"required,oneof=credit_card debit_card pix bolbradesco"`
	Descricao       string  `json:"descricao" validate:"max=500"`
	NumeroParcelas  int     `json:"numero_parcelas" validate:"min=1,max=12"`
	TokenCartao     string  `json:"token_cartao"`
	ChavePIX        string  `json:"chave_pix"`
}

// Validate - valida os campos da requisição
func (r *CreatePagamentoRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ClienteID, validation.Required, validation.Min(1)),
		validation.Field(&r.EmpresaID, validation.Required, validation.Min(1)),
		validation.Field(&r.Valor, validation.Required, validation.Min(0.01)),
		validation.Field(&r.Moeda, validation.Required, validation.In("BRL", "USD", "EUR", "ARS", "CLP", "COP", "MXN", "PEN", "UYU")),
		validation.Field(&r.MetodoPagamento, validation.Required, validation.In("credit_card", "debit_card", "pix", "bolbradesco")),
		validation.Field(&r.Descricao, validation.Length(0, 500)),
		validation.Field(&r.NumeroParcelas, validation.Min(1), validation.Max(12)),
	)
}

// UpdatePagamentoRequest - representa a requisição para atualizar um pagamento
type UpdatePagamentoRequest struct {
	Status       string `json:"status" validate:"omitempty,oneof=pending approved authorized in_process in_mediation rejected cancelled refunded charged_back"`
	StatusDetail string `json:"status_detail" validate:"max=255"`
	TokenCartao  string `json:"token_cartao"`
	ChavePIX     string `json:"chave_pix"`
	QRCode       string `json:"qr_code"`
}

// Validate - valida os campos da requisição
func (r *UpdatePagamentoRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Status, validation.In("pending", "approved", "authorized", "in_process", "in_mediation", "rejected", "cancelled", "refunded", "charged_back")),
		validation.Field(&r.StatusDetail, validation.Length(0, 255)),
	)
}

// GetPagamentoRequest - representa a requisição para buscar um pagamento
type GetPagamentoRequest struct {
	ID int `json:"id" validate:"required,min=1"`
}

// Validate - valida os campos da requisição
func (r *GetPagamentoRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ID, validation.Required, validation.Min(1)),
	)
}

// ListPagamentoRequest - representa a requisição para listar pagamentos
type ListPagamentoRequest struct {
	ClienteID       int    `json:"cliente_id" validate:"omitempty,min=1"`
	EmpresaID       int    `json:"empresa_id" validate:"omitempty,min=1"`
	Status          string `json:"status" validate:"omitempty,oneof=pending approved authorized in_process in_mediation rejected cancelled refunded charged_back"`
	MetodoPagamento string `json:"metodo_pagamento" validate:"omitempty,oneof=credit_card debit_card pix bolbradesco"`
	Page            int    `json:"page" validate:"min=1"`
	Limit           int    `json:"limit" validate:"min=1,max=100"`
}

// Validate - valida os campos da requisição
func (r *ListPagamentoRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ClienteID, validation.Min(1)),
		validation.Field(&r.EmpresaID, validation.Min(1)),
		validation.Field(&r.Status, validation.In("pending", "approved", "authorized", "in_process", "in_mediation", "rejected", "cancelled", "refunded", "charged_back")),
		validation.Field(&r.MetodoPagamento, validation.In("credit_card", "debit_card", "pix", "bolbradesco")),
		validation.Field(&r.Page, validation.Min(1)),
		validation.Field(&r.Limit, validation.Min(1), validation.Max(100)),
	)
}
