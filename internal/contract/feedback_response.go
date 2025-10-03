package contract

import "time"

// FeedbackResponse - representa a resposta de um feedback
type FeedbackResponse struct {
	ID                 int       `json:"id"`
	ClienteID          int       `json:"cliente_id"`
	EmpresaID          int       `json:"empresa_id"`
	ReservaID          int       `json:"reserva_id"`
	Nota               int       `json:"nota"`
	Comentario         string    `json:"comentario"`
	Status             string    `json:"status"`
	MomentoCriacao     time.Time `json:"momento_criacao"`
	MomentoAtualizacao time.Time `json:"momento_atualizacao"`
	StatusDisplay      string    `json:"status_display"`
	RatingDisplay      string    `json:"rating_display"`
}

// ListFeedbackResponse - representa a resposta de uma lista de feedbacks
type ListFeedbackResponse struct {
	Feedbacks []FeedbackResponse `json:"feedbacks"`
	Total     int                `json:"total"`
	Page      int                `json:"page"`
	Limit     int                `json:"limit"`
	Pages     int                `json:"pages"`
}

// CreateFeedbackResponse - representa a resposta da criação de um feedback
type CreateFeedbackResponse struct {
	Feedback FeedbackResponse `json:"feedback"`
	Message  string           `json:"message"`
}

// UpdateFeedbackResponse - representa a resposta da atualização de um feedback
type UpdateFeedbackResponse struct {
	Feedback FeedbackResponse `json:"feedback"`
	Message  string           `json:"message"`
}
