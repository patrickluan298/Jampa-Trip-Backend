package contract

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
)

// CreateReservaRequest - representa a requisição para criar uma reserva
type CreateReservaRequest struct {
	TourID                 int       `json:"tour_id" validate:"required,min=1"`
	ClienteID              int       `json:"cliente_id" validate:"required,min=1"`
	PagamentoID            int       `json:"pagamento_id" validate:"omitempty,min=1"`
	DataPasseioSelecionada time.Time `json:"data_passeio_selecionada" validate:"required"`
	QuantidadePessoas      int       `json:"quantidade_pessoas" validate:"required,min=1,max=50"`
	Observacoes            string    `json:"observacoes" validate:"max=1000"`
}

// Validate - valida os campos da requisição
func (r *CreateReservaRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.TourID, validation.Required, validation.Min(1)),
		validation.Field(&r.ClienteID, validation.Required, validation.Min(1)),
		validation.Field(&r.PagamentoID, validation.Min(1)),
		validation.Field(&r.DataPasseioSelecionada, validation.Required),
		validation.Field(&r.QuantidadePessoas, validation.Required, validation.Min(1), validation.Max(50)),
		validation.Field(&r.Observacoes, validation.Length(0, 1000)),
	)
}

// UpdateReservaRequest - representa a requisição para atualizar uma reserva
type UpdateReservaRequest struct {
	Status                 string    `json:"status" validate:"omitempty,oneof=pendente aguardando_pagamento confirmada cancelada concluida"`
	DataPasseioSelecionada time.Time `json:"data_passeio_selecionada" validate:"omitempty"`
	QuantidadePessoas      int       `json:"quantidade_pessoas" validate:"omitempty,min=1,max=50"`
	Observacoes            string    `json:"observacoes" validate:"max=1000"`
}

// Validate - valida os campos da requisição
func (r *UpdateReservaRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Status, validation.In("pendente", "aguardando_pagamento", "confirmada", "cancelada", "concluida")),
		validation.Field(&r.QuantidadePessoas, validation.Min(1), validation.Max(50)),
		validation.Field(&r.Observacoes, validation.Length(0, 1000)),
	)
}

// GetReservaRequest - representa a requisição para buscar uma reserva
type GetReservaRequest struct {
	ID int `json:"id" validate:"required,min=1"`
}

// Validate - valida os campos da requisição
func (r *GetReservaRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ID, validation.Required, validation.Min(1)),
	)
}

// ListReservaRequest - representa a requisição para listar reservas
type ListReservaRequest struct {
	TourID    int    `json:"tour_id" validate:"omitempty,min=1"`
	ClienteID int    `json:"cliente_id" validate:"omitempty,min=1"`
	CompanyID int    `json:"company_id" validate:"omitempty,min=1"`
	Status    string `json:"status" validate:"omitempty,oneof=pendente aguardando_pagamento confirmada cancelada concluida"`
	Page      int    `json:"page" validate:"min=1"`
	Limit     int    `json:"limit" validate:"min=1,max=100"`
}

// Validate - valida os campos da requisição
func (r *ListReservaRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.TourID, validation.Min(1)),
		validation.Field(&r.ClienteID, validation.Min(1)),
		validation.Field(&r.CompanyID, validation.Min(1)),
		validation.Field(&r.Status, validation.In("pendente", "aguardando_pagamento", "confirmada", "cancelada", "concluida")),
		validation.Field(&r.Page, validation.Min(1)),
		validation.Field(&r.Limit, validation.Min(1), validation.Max(100)),
	)
}

// CancelarReservaRequest - representa a requisição para cancelar uma reserva
type CancelarReservaRequest struct {
	ID int `json:"id" validate:"required,min=1"`
}

// Validate - valida os campos da requisição
func (r *CancelarReservaRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ID, validation.Required, validation.Min(1)),
	)
}
