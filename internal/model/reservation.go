package model

import "time"

// Reserva - representa a entidade de reserva
type Reserva struct {
	ID                     int        `gorm:"column:id;primaryKey;autoIncrement"`
	TourID                 int        `gorm:"column:tour_id;not null"`
	ClienteID              int        `gorm:"column:cliente_id;not null"`
	PagamentoID            int        `gorm:"column:pagamento_id"`
	Status                 string     `gorm:"column:status;not null;default:'pendente'"`
	DataReserva            time.Time  `gorm:"column:data_reserva;not null"`
	DataPasseioSelecionada time.Time  `gorm:"column:data_passeio_selecionada;not null;type:date"`
	QuantidadePessoas      int        `gorm:"column:quantidade_pessoas;not null;default:1"`
	ValorTotal             float64    `gorm:"column:valor_total;not null;type:decimal(10,2)"`
	Observacoes            string     `gorm:"column:observacoes"`
	MomentoCriacao         time.Time  `gorm:"column:momento_criacao;not null;default:CURRENT_TIMESTAMP"`
	MomentoAtualizacao     time.Time  `gorm:"column:momento_atualizacao;not null;default:CURRENT_TIMESTAMP"`
	MomentoCancelamento    *time.Time `gorm:"column:momento_cancelamento"`

	// Relacionamentos
	Tour      Tour      `gorm:"foreignKey:TourID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Cliente   Client    `gorm:"foreignKey:ClienteID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Pagamento Pagamento `gorm:"foreignKey:PagamentoID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

// TableName - especifica o nome da tabela no banco de dados
func (Reserva) TableName() string {
	return "reservations"
}

// StatusReserva - define os possíveis status de uma reserva
type StatusReserva string

const (
	StatusReservaPendente            StatusReserva = "pendente"
	StatusReservaAguardandoPagamento StatusReserva = "aguardando_pagamento"
	StatusReservaConfirmada          StatusReserva = "confirmada"
	StatusReservaCancelada           StatusReserva = "cancelada"
	StatusReservaConcluida           StatusReserva = "concluida"
)

// IsPending - verifica se a reserva está pendente
func (r *Reserva) IsPending() bool {
	return r.Status == string(StatusReservaPendente)
}

// IsAwaitingPayment - verifica se a reserva está aguardando pagamento
func (r *Reserva) IsAwaitingPayment() bool {
	return r.Status == string(StatusReservaAguardandoPagamento)
}

// IsConfirmed - verifica se a reserva está confirmada
func (r *Reserva) IsConfirmed() bool {
	return r.Status == string(StatusReservaConfirmada)
}

// IsCancelled - verifica se a reserva está cancelada
func (r *Reserva) IsCancelled() bool {
	return r.Status == string(StatusReservaCancelada)
}

// IsCompleted - verifica se a reserva está concluída
func (r *Reserva) IsCompleted() bool {
	return r.Status == string(StatusReservaConcluida)
}

// CanBeCancelled - verifica se a reserva pode ser cancelada baseado no status
func (r *Reserva) CanBeCancelled() bool {
	return r.Status == string(StatusReservaPendente) ||
		r.Status == string(StatusReservaAguardandoPagamento) ||
		r.Status == string(StatusReservaConfirmada)
}

// GetStatusDisplay - retorna o status em formato legível
func (r *Reserva) GetStatusDisplay() string {
	switch r.Status {
	case string(StatusReservaPendente):
		return "Pendente"
	case string(StatusReservaAguardandoPagamento):
		return "Aguardando Pagamento"
	case string(StatusReservaConfirmada):
		return "Confirmada"
	case string(StatusReservaCancelada):
		return "Cancelada"
	case string(StatusReservaConcluida):
		return "Concluída"
	default:
		return "Desconhecido"
	}
}
