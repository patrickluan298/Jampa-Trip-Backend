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

// AutorizarCartaoRequest - representa a requisição para autorizar um pagamento com cartão
type AutorizarCartaoRequest struct {
	ClienteID         int          `json:"cliente_id" validate:"required,min=1"`
	EmpresaID         int          `json:"empresa_id" validate:"required,min=1"`
	Token             string       `json:"token" validate:"required,min=1"`
	TransactionAmount float64      `json:"transaction_amount" validate:"required,min=0.01"`
	Installments      int          `json:"installments" validate:"required,min=1,max=12"`
	PaymentMethodID   string       `json:"payment_method_id" validate:"required"`
	IssuerID          string       `json:"issuer_id"`
	Description       string       `json:"description" validate:"max=500"`
	Payer             PayerRequest `json:"payer" validate:"required"`
	ExternalReference string       `json:"external_reference"`
}

// PayerRequest - representa os dados do pagador
type PayerRequest struct {
	Email          string                `json:"email" validate:"required,email"`
	Identification IdentificationRequest `json:"identification"`
	FirstName      string                `json:"first_name"`
	LastName       string                `json:"last_name"`
}

// IdentificationRequest - representa a identificação do pagador
type IdentificationRequest struct {
	Type   string `json:"type" validate:"required,oneof=CPF CNPJ"`
	Number string `json:"number" validate:"required"`
}

// Validate - valida os campos da requisição
func (r *AutorizarCartaoRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ClienteID, validation.Required, validation.Min(1)),
		validation.Field(&r.EmpresaID, validation.Required, validation.Min(1)),
		validation.Field(&r.Token, validation.Required, validation.Length(1, 500)),
		validation.Field(&r.TransactionAmount, validation.Required, validation.Min(0.01)),
		validation.Field(&r.Installments, validation.Required, validation.Min(1), validation.Max(12)),
		validation.Field(&r.PaymentMethodID, validation.Required, validation.In("visa", "master", "amex", "elo", "hipercard", "cabal", "naranja", "tarshop")),
		validation.Field(&r.Description, validation.Length(0, 500)),
		validation.Field(&r.Payer.Email, validation.Required),
		validation.Field(&r.Payer.Identification.Type, validation.In("CPF", "CNPJ")),
	)
}

// CapturarPagamentoRequest - representa a requisição para capturar um pagamento autorizado
type CapturarPagamentoRequest struct {
	PaymentID         int64    `json:"payment_id" validate:"required,min=1"`
	TransactionAmount *float64 `json:"transaction_amount,omitempty"`
}

// Validate - valida os campos da requisição
func (r *CapturarPagamentoRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.PaymentID, validation.Required, validation.Min(1)),
	)
}

// CancelarPagamentoRequest - representa a requisição para cancelar um pagamento
type CancelarPagamentoRequest struct {
	PaymentID int64 `json:"payment_id" validate:"required,min=1"`
}

// Validate - valida os campos da requisição
func (r *CancelarPagamentoRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.PaymentID, validation.Required, validation.Min(1)),
	)
}

// ReembolsarPagamentoRequest - representa a requisição para reembolsar um pagamento
type ReembolsarPagamentoRequest struct {
	PaymentID int64    `json:"payment_id" validate:"required,min=1"`
	Amount    *float64 `json:"amount,omitempty"`
}

// Validate - valida os campos da requisição
func (r *ReembolsarPagamentoRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.PaymentID, validation.Required, validation.Min(1)),
	)
}
