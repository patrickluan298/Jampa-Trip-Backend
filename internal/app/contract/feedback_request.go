package contract

import validation "github.com/go-ozzo/ozzo-validation"

// CreateFeedbackRequest - representa a requisição para criar um feedback
type CreateFeedbackRequest struct {
	ClienteID  int    `json:"cliente_id" validate:"required,min=1"`
	EmpresaID  int    `json:"empresa_id" validate:"required,min=1"`
	ReservaID  int    `json:"reserva_id" validate:"omitempty,min=1"`
	Nota       int    `json:"nota" validate:"required,min=1,max=5"`
	Comentario string `json:"comentario" validate:"max=1000"`
}

// Validate - valida os campos da requisição
func (r *CreateFeedbackRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ClienteID, validation.Required, validation.Min(1)),
		validation.Field(&r.EmpresaID, validation.Required, validation.Min(1)),
		validation.Field(&r.ReservaID, validation.Min(1)),
		validation.Field(&r.Nota, validation.Required, validation.Min(1), validation.Max(5)),
		validation.Field(&r.Comentario, validation.Length(0, 1000)),
	)
}

// UpdateFeedbackRequest - representa a requisição para atualizar um feedback
type UpdateFeedbackRequest struct {
	Nota       int    `json:"nota" validate:"omitempty,min=1,max=5"`
	Comentario string `json:"comentario" validate:"max=1000"`
	Status     string `json:"status" validate:"omitempty,oneof=ativo inativo moderado"`
}

// Validate - valida os campos da requisição
func (r *UpdateFeedbackRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Nota, validation.Min(1), validation.Max(5)),
		validation.Field(&r.Comentario, validation.Length(0, 1000)),
		validation.Field(&r.Status, validation.In("ativo", "inativo", "moderado")),
	)
}

// GetFeedbackRequest - representa a requisição para buscar um feedback
type GetFeedbackRequest struct {
	ID int `json:"id" validate:"required,min=1"`
}

// Validate - valida os campos da requisição
func (r *GetFeedbackRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ID, validation.Required, validation.Min(1)),
	)
}

// ListFeedbackRequest - representa a requisição para listar feedbacks
type ListFeedbackRequest struct {
	ClienteID int    `json:"cliente_id" validate:"omitempty,min=1"`
	EmpresaID int    `json:"empresa_id" validate:"omitempty,min=1"`
	Status    string `json:"status" validate:"omitempty,oneof=ativo inativo moderado"`
	Nota      int    `json:"nota" validate:"omitempty,min=1,max=5"`
	Page      int    `json:"page" validate:"min=1"`
	Limit     int    `json:"limit" validate:"min=1,max=100"`
}

// Validate - valida os campos da requisição
func (r *ListFeedbackRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ClienteID, validation.Min(1)),
		validation.Field(&r.EmpresaID, validation.Min(1)),
		validation.Field(&r.Status, validation.In("ativo", "inativo", "moderado")),
		validation.Field(&r.Nota, validation.Min(1), validation.Max(5)),
		validation.Field(&r.Page, validation.Min(1)),
		validation.Field(&r.Limit, validation.Min(1), validation.Max(100)),
	)
}
