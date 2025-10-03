package model

import "time"

// Feedback - representa a entidade de feedback
type Feedback struct {
	ID                 int       `gorm:"column:id;primaryKey;autoIncrement"`
	ClienteID          int       `gorm:"column:cliente_id;not null"`
	EmpresaID          int       `gorm:"column:empresa_id;not null"`
	ReservaID          int       `gorm:"column:reserva_id"`
	Nota               int       `gorm:"column:nota;not null;check:nota >= 1 AND nota <= 5"`
	Comentario         string    `gorm:"column:comentario"`
	Status             string    `gorm:"column:status;not null;default:'ativo'"`
	MomentoCriacao     time.Time `gorm:"column:momento_criacao;not null;default:CURRENT_TIMESTAMP"`
	MomentoAtualizacao time.Time `gorm:"column:momento_atualizacao;not null;default:CURRENT_TIMESTAMP"`

	// Relacionamentos
	Cliente Cliente `gorm:"foreignKey:ClienteID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Empresa Empresa `gorm:"foreignKey:EmpresaID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Reserva Reserva `gorm:"foreignKey:ReservaID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

// TableName - especifica o nome da tabela no banco de dados
func (Feedback) TableName() string {
	return "feedbacks"
}

// StatusFeedback - define os possíveis status de um feedback
type StatusFeedback string

const (
	StatusFeedbackAtivo    StatusFeedback = "ativo"
	StatusFeedbackInativo  StatusFeedback = "inativo"
	StatusFeedbackModerado StatusFeedback = "moderado"
)

// Métodos de validação para Feedback
func (f *Feedback) IsValid() bool {
	return f.ClienteID > 0 && f.EmpresaID > 0 && f.Nota >= 1 && f.Nota <= 5
}

func (f *Feedback) IsActive() bool {
	return f.Status == string(StatusFeedbackAtivo)
}

func (f *Feedback) IsInactive() bool {
	return f.Status == string(StatusFeedbackInativo)
}

func (f *Feedback) IsModerated() bool {
	return f.Status == string(StatusFeedbackModerado)
}

func (f *Feedback) UpdateStatus(status StatusFeedback) {
	f.Status = string(status)
	f.MomentoAtualizacao = time.Now()
}

func (f *Feedback) GetStatusDisplay() string {
	switch f.Status {
	case string(StatusFeedbackAtivo):
		return "Ativo"
	case string(StatusFeedbackInativo):
		return "Inativo"
	case string(StatusFeedbackModerado):
		return "Moderado"
	default:
		return "Desconhecido"
	}
}

func (f *Feedback) GetRatingDisplay() string {
	switch f.Nota {
	case 1:
		return "Muito Ruim"
	case 2:
		return "Ruim"
	case 3:
		return "Regular"
	case 4:
		return "Bom"
	case 5:
		return "Excelente"
	default:
		return "Desconhecido"
	}
}

// Funções de validação
func IsValidFeedbackStatus(status StatusFeedback) bool {
	switch status {
	case StatusFeedbackAtivo, StatusFeedbackInativo, StatusFeedbackModerado:
		return true
	default:
		return false
	}
}

func IsValidRating(rating int) bool {
	return rating >= 1 && rating <= 5
}
