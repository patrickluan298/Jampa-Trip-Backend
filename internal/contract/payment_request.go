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

// CreateCreditCardPaymentRequest - representa a requisição para criar pagamento com cartão de crédito
type CreateCreditCardPaymentRequest struct {
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
	Capture           bool         `json:"capture"`
}

// Validate - valida os campos da requisição
func (r *CreateCreditCardPaymentRequest) Validate() error {
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

// CreateDebitCardPaymentRequest - representa a requisição para criar pagamento com cartão de débito
type CreateDebitCardPaymentRequest struct {
	ClienteID         int          `json:"cliente_id" validate:"required,min=1"`
	EmpresaID         int          `json:"empresa_id" validate:"required,min=1"`
	Token             string       `json:"token" validate:"required,min=1"`
	TransactionAmount float64      `json:"transaction_amount" validate:"required,min=0.01"`
	PaymentMethodID   string       `json:"payment_method_id" validate:"required"`
	IssuerID          string       `json:"issuer_id"`
	Description       string       `json:"description" validate:"max=500"`
	Payer             PayerRequest `json:"payer" validate:"required"`
	ExternalReference string       `json:"external_reference"`
}

// Validate - valida os campos da requisição
func (r *CreateDebitCardPaymentRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ClienteID, validation.Required, validation.Min(1)),
		validation.Field(&r.EmpresaID, validation.Required, validation.Min(1)),
		validation.Field(&r.Token, validation.Required, validation.Length(1, 500)),
		validation.Field(&r.TransactionAmount, validation.Required, validation.Min(0.01)),
		validation.Field(&r.PaymentMethodID, validation.Required, validation.In("visa", "master", "amex", "elo", "hipercard", "cabal", "naranja", "tarshop")),
		validation.Field(&r.Description, validation.Length(0, 500)),
		validation.Field(&r.Payer.Email, validation.Required),
		validation.Field(&r.Payer.Identification.Type, validation.In("CPF", "CNPJ")),
	)
}

// CreatePIXPaymentRequest - representa a requisição para criar pagamento com PIX
type CreatePIXPaymentRequest struct {
	ClienteID         int          `json:"cliente_id" validate:"required,min=1"`
	EmpresaID         int          `json:"empresa_id" validate:"required,min=1"`
	TransactionAmount float64      `json:"transaction_amount" validate:"required,min=0.01"`
	Description       string       `json:"description" validate:"max=500"`
	Payer             PayerRequest `json:"payer" validate:"required"`
	ExternalReference string       `json:"external_reference"`
}

// Validate - valida os campos da requisição
func (r *CreatePIXPaymentRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ClienteID, validation.Required, validation.Min(1)),
		validation.Field(&r.EmpresaID, validation.Required, validation.Min(1)),
		validation.Field(&r.TransactionAmount, validation.Required, validation.Min(0.01)),
		validation.Field(&r.Description, validation.Length(0, 500)),
		validation.Field(&r.Payer.Email, validation.Required),
	)
}

// ListPaymentsRequest - representa a requisição para buscar pagamentos
type ListPaymentsRequest struct {
	ExternalReference string `json:"external_reference"`
	Status            string `json:"status"`
	PaymentMethod     string `json:"payment_method"`
	DateCreatedFrom   string `json:"date_created_from"`
	DateCreatedTo     string `json:"date_created_to"`
	Offset            int    `json:"offset" validate:"min=0"`
	Limit             int    `json:"limit" validate:"min=1,max=100"`
}

// Validate - valida os campos da requisição
func (r *ListPaymentsRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Status, validation.In("pending", "approved", "authorized", "in_process", "in_mediation", "rejected", "cancelled", "refunded", "charged_back")),
		validation.Field(&r.PaymentMethod, validation.In("credit_card", "debit_card", "pix", "bolbradesco")),
		validation.Field(&r.Offset, validation.Min(0)),
		validation.Field(&r.Limit, validation.Min(1), validation.Max(100)),
	)
}

// ObterPagamentoRequest - representa a requisição para obter um pagamento por ID
type ObterPagamentoRequest struct {
	ID int64 `json:"id" validate:"required,min=1"`
}

// Validate - valida os campos da requisição
func (r *ObterPagamentoRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ID, validation.Required, validation.Min(1)),
	)
}

// UpdatePaymentRequest - representa a requisição para atualizar um pagamento
type UpdatePaymentRequest struct {
	ID              int64             `json:"id" validate:"required,min=1"`
	Description     string            `json:"description,omitempty"`
	Metadata        map[string]string `json:"metadata,omitempty"`
	NotificationURL string            `json:"notification_url,omitempty"`
}

// Validate - valida os campos da requisição
func (r *UpdatePaymentRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ID, validation.Required, validation.Min(1)),
		validation.Field(&r.Description, validation.Length(0, 500)),
	)
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
