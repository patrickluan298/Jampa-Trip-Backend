package contract

import "time"

// ReservaResponse - representa a resposta de uma reserva
type ReservaResponse struct {
	ID                  int        `json:"id"`
	ClienteID           int        `json:"cliente_id"`
	EmpresaID           int        `json:"empresa_id"`
	PagamentoID         int        `json:"pagamento_id"`
	Status              string     `json:"status"`
	DataReserva         time.Time  `json:"data_reserva"`
	DataPasseio         time.Time  `json:"data_passeio"`
	QuantidadePessoas   int        `json:"quantidade_pessoas"`
	ValorTotal          float64    `json:"valor_total"`
	Observacoes         string     `json:"observacoes"`
	MomentoCriacao      time.Time  `json:"momento_criacao"`
	MomentoAtualizacao  time.Time  `json:"momento_atualizacao"`
	MomentoCancelamento *time.Time `json:"momento_cancelamento"`
	StatusDisplay       string     `json:"status_display"`
}

// ListReservaResponse - representa a resposta de uma lista de reservas
type ListReservaResponse struct {
	Reservas []ReservaResponse `json:"reservas"`
	Total    int               `json:"total"`
	Page     int               `json:"page"`
	Limit    int               `json:"limit"`
	Pages    int               `json:"pages"`
}

// CreateReservaResponse - representa a resposta da criação de uma reserva
type CreateReservaResponse struct {
	Reserva ReservaResponse `json:"reserva"`
	Message string          `json:"message"`
}

// UpdateReservaResponse - representa a resposta da atualização de uma reserva
type UpdateReservaResponse struct {
	Reserva ReservaResponse `json:"reserva"`
	Message string          `json:"message"`
}

// CancelarReservaResponse - representa a resposta do cancelamento de uma reserva
type CancelarReservaResponse struct {
	Reserva ReservaResponse `json:"reserva"`
	Message string          `json:"message"`
}
